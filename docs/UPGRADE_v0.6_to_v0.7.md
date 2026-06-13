# Upgrade Guide: v0.6 to v0.7

This document outlines the breaking changes and upgrade considerations when migrating from sacloud-otel-collector v0.6.x (based on OpenTelemetry Collector v0.142.0) to v0.7.x (based on OpenTelemetry Collector v0.154.0).

## OpenTelemetry Collector Version Update

- **v0.6.x**: Based on OpenTelemetry Collector v0.142.0
- **v0.7.x**: Based on OpenTelemetry Collector v0.154.0

## Why This Upgrade

This upgrade is primarily driven by security fixes. Several transitive dependencies pinned by Collector v0.142.0 (notably `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/sdk`, and `github.com/prometheus/prometheus`) had open security advisories that cannot be patched without moving the Collector forward. v0.154.0 ships `otel-sdk` >= 1.43.0 and `prometheus` >= 0.311.3, resolving those advisories.

## Impact Summary

The components used by the default `config.yaml` (`otlp` receiver, `batch` processor, `debug` exporter, `health_check` extension) are **not affected** by any breaking change in this range. If you only use those, your existing configuration continues to validate and run unchanged.

The breaking changes below matter only if your configuration uses the affected components. The two most important classes of change are:

1. **Silent behavior change** — `filter`/`transform` processors now default to `error_mode: ignore`.
2. **Config-invalidating removals** — `prometheus` receiver and `kafka` receiver removed several config fields; configs that still set them will be rejected by `validate`.

## Breaking Changes by Component

### Core Collector Changes

#### 1. Go Version Requirement

**v0.146.0 Changes:**
- Minimum Go version increased to 1.25
- **Impact**: Only affects users building from source

#### 2. Internal telemetry: default metrics address moved to localhost

**v0.153.0 Changes:**
- `telemetry.UseLocalHostAsDefaultMetricsAddress` feature gate stabilized ([#15342](https://github.com/open-telemetry/opentelemetry-collector/pull/15342))
- **Impact**: The Collector's own internal metrics endpoint now defaults to `localhost:8888` instead of `0.0.0.0:8888`. If you scrape the Collector's self-metrics from another host or container, the default no longer listens on all interfaces.
- **Action Required (conditional)**: To restore remote scraping, set the address explicitly under `service::telemetry::metrics`.

#### 3. Internal telemetry: constant labels removed

**v0.149.0 Changes:**
- `service_name`, `service_instance_id`, and `service_version` removed as constant labels on every internal metric datapoint ([#14811](https://github.com/open-telemetry/opentelemetry-collector/pull/14811))
- **Impact**: These values remain available via `target_info`. Only affects dashboards/alerts that group or filter the Collector's **own** internal metrics by those labels. Telemetry passing through the Collector is unaffected.

### Components with Breaking Changes

#### 1. transformprocessor / filterprocessor

These two share the OTTL engine and have the most significant behavioral changes in this range.

**v0.153.0 Changes:**
- `processor/filter`: default `error_mode` changed from `propagate` to `ignore` ([#47232](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/47232))
- `processor/transform`: default top-level `error_mode` changed from `propagate` to `ignore` ([#48415](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/48415))
- **Impact**: Errors raised by filter/transform statements are now **silently ignored** by default instead of failing the pipeline.
- **Action Required**: If you relied on errors propagating, set `error_mode: propagate` explicitly in the processor config, or disable the gate (`--feature-gates=-processor.filter.defaultErrorModeIgnore` / `--feature-gates=-processor.transform.defaultErrorModeIgnore`).

**v0.154.0 Changes:**
- `pkg/ottl`: OTTL datapoint context setters now return an error when used on an incompatible data point type ([#48384](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/48384))
- **Impact**: A statement such as `set(explicit_bounds, [1.0])` against a `NumberDataPoint` previously no-op'd silently; it now errors. Combined with the new default `error_mode: ignore`, such a statement silently does nothing rather than erroring the pipeline.
- **Action Required**: Audit transform/filter OTTL statements that set datapoint fields (`value_double`/`value_int` for Number, `explicit_bounds`/`bucket_counts` for Histogram, `scale`/`zero_count`/`positive*`/`negative*` for ExponentialHistogram, `quantile_values` for Summary).

#### 2. prometheusreceiver

Configurations that still set the removed options will be **rejected by validation**.

**v0.143.0 Changes:**
- Removed deprecated `use_start_time_metric` and `start_time_metric_regex` options ([#44180](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/44180))
- **Action Required**: Remove these fields. If you relied on start-time adjustment, use the `metricstarttime` processor.

**v0.149.0 Changes:**
- Removed the `report_extra_scrape_metrics` option and the obsolete extra-scrape-metric feature gates ([#44181](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/44181))
- **Action Required**: Remove `report_extra_scrape_metrics` from the receiver config.

#### 3. kafkareceiver

The Sarama client is fully removed; franz-go is now the only implementation.

**v0.144.0 Changes:**
- Sarama consumer implementation removed; `default_fetch_size` option removed; `receiver.kafkareceiver.UseFranzGo` feature gate removed ([#44564](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/44564))
- **Action Required**: Remove `default_fetch_size` (use `max_fetch_size`) and drop any `--feature-gates` reference to `receiver.kafkareceiver.UseFranzGo`.

**v0.147.0 Changes:**
- Removed deprecated `topic` and `exclude_topic` fields ([#46232](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/46232))
- **Action Required**: Use the list-based `topics` / `exclude_topics` configuration instead.
- **Example**:
  ```yaml
  # Old configuration (v0.6.x)
  receivers:
    kafka:
      logs:
        topic: my-topic
        exclude_topic: excluded-topic

  # New configuration (v0.7.x)
  receivers:
    kafka:
      logs:
        topics:
          - my-topic
        exclude_topics:
          - excluded-topic
  ```

#### 4. tailsamplingprocessor

**v0.144.0 Changes:**
- Deprecated invert decisions disabled by default ([#44132](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/44132))
- **Impact**: `invert_match: true` on string/numeric/boolean attribute policies no longer inverts by default.
- **Action Required**: Migrate to `drop` policies for explicit non-sampling, or temporarily disable the gate (`--feature-gates=-processor.tailsamplingprocessor.disableinvertdecisions`).

#### 5. resourcedetectionprocessor

**v0.146.0 Changes:**
- `processor.resourcedetection.propagateerrors` feature gate promoted to Stable and always enabled ([#44609](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/44609))

**v0.150.0 Changes:**
- The `processor.resourcedetection.propagateerrors` feature gate was removed entirely ([#45853](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/45853))
- **Impact**: Resource-detector errors now always propagate and can fail the pipeline; there is **no longer an opt-out**. In environments where a configured detector can fail (e.g. an unreachable cloud-metadata endpoint), startup may now error where it previously continued.
- **Action Required**: Review your `detectors:` list to ensure each configured detector can succeed in the target environment.

#### 6. hostmetricsreceiver

**v0.146.0 Changes:**
- `process.context_switches` now counts context switches for all threads, not just the lead thread ([#36804](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/36804))
- **Impact**: No config change, but the metric's values may change drastically. Re-baseline any alerts that use this metric.

#### 7. elasticsearchexporter (ECS mapping mode only)

A series of releases removed ECS-mode enrichment fields. These only matter when using `mapping::mode: ecs`:
- **v0.144.0**: removed ECS span enrichment for `span.action`, `span.message.queue.name`, `transaction.message.queue.name`
- **v0.146.0**: removed ECS log enrichment for `agent.name`, `agent.version`
- **v0.149.0**: removed `host.os.type` encoding in ECS mode (use `processor/elasticapmprocessor` for that enrichment)
- **Action Required (conditional)**: If you depend on these fields in ECS mode, add the external `elasticapmprocessor` or accept the dropped fields.

### Components with No Breaking Changes

The following bundled components had no breaking changes in the v0.143.0–v0.154.0 range:
`otlpreceiver`, `otlpexporter`, `otlphttpexporter`, `debugexporter`, `batchprocessor`, `memorylimiterprocessor`, `awss3exporter`, `fileexporter`, `prometheusremotewriteexporter`, `healthcheckextension`, `filestorage`, `attributesprocessor`, `deltatocumulativeprocessor`, `groupbyattrsprocessor`, `resourceprocessor`, `filelogreceiver`, `fluentforwardreceiver`, `journaldreceiver`.

> Note: `mackerelotlpexporter` is a third-party component (`github.com/mackerelio/opentelemetry-collector-mackerel`) with its own versioning and was not covered by the Collector changelogs above. Review its release notes separately if you use it.

## Configuration Validation

After updating configurations, validate them using:

```bash
./sacloud-otel-collector validate --config your-config.yaml
```

The `prometheus` and `kafka` receiver field removals are the changes most likely to make an existing `config.yaml` fail validation outright. The `filter`/`transform` `error_mode` change is the one most likely to silently change behavior while still validating.

## References

### OpenTelemetry Collector Releases
- [v0.143.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.143.0)
- [v0.144.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.144.0)
- [v0.145.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.145.0)
- [v0.146.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.146.0)
- [v0.147.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.147.0)
- [v0.148.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.148.0)
- [v0.149.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.149.0)
- [v0.150.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.150.0)
- [v0.151.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.151.0)
- [v0.152.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.152.0)
- [v0.153.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.153.0)
- [v0.154.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.154.0)

### OpenTelemetry Collector Contrib Releases
- [v0.143.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.143.0)
- [v0.144.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.144.0)
- [v0.145.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.145.0)
- [v0.146.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.146.0)
- [v0.147.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.147.0)
- [v0.148.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.148.0)
- [v0.149.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.149.0)
- [v0.150.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.150.0)
- [v0.151.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.151.0)
- [v0.152.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.152.0)
- [v0.153.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.153.0)
- [v0.154.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.154.0)

## Support

If you encounter issues during the upgrade:

1. Check component-specific documentation at [OpenTelemetry Docs](https://opentelemetry.io/docs/)
2. Review the full changelog for each version
3. Report issues at the [sacloud-otel-collector repository](https://github.com/sacloud/sacloud-otel-collector/issues)
