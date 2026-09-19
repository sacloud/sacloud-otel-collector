# Upgrade Guide: v0.7 to v0.8

This document outlines the breaking changes and upgrade considerations when migrating from sacloud-otel-collector v0.7.x (based on OpenTelemetry Collector v0.154.0) to v0.8.x (based on OpenTelemetry Collector v0.161.0).

## OpenTelemetry Collector Version Update

- **v0.7.x**: Based on OpenTelemetry Collector v0.154.0
- **v0.8.x**: Based on OpenTelemetry Collector v0.161.0

## Why This Upgrade

This upgrade is primarily driven by security fixes. Several transitive dependencies pinned by Collector v0.154.0 had open security advisories that cannot be patched without moving the Collector forward:

- `google.golang.org/grpc` — HTTP/2 DATA frame fragmentation OOM ([GHSA-vp52-pcj8-j9qc](https://github.com/advisories/GHSA-vp52-pcj8-j9qc)), xDS server crash ([GHSA-2v4p-qf9q-27wj](https://github.com/advisories/GHSA-2v4p-qf9q-27wj)), xDS RBAC bypass ([GHSA-qc2q-p7wx-3px3](https://github.com/advisories/GHSA-qc2q-p7wx-3px3))
- `github.com/apache/thrift` — infinite loop ([GHSA-8wv5-x4w7-5gww](https://github.com/advisories/GHSA-8wv5-x4w7-5gww))
- `go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc` — TLS certificates from environment variables ignored ([GHSA-w34q-cm8f-9c5x](https://github.com/advisories/GHSA-w34q-cm8f-9c5x))
- `go.opentelemetry.io/otel/sdk` and OTLP trace exporters — endpoint URLs leaked in info logs ([GHSA-8wmf-6v46-5gfg](https://github.com/advisories/GHSA-8wmf-6v46-5gfg))
- `github.com/cilium/ebpf` — integer overflow ([GHSA-xhgw-qwwf-pg32](https://github.com/advisories/GHSA-xhgw-qwwf-pg32))

v0.161.0 ships `grpc` v1.83.2, `thrift` v0.24.0, `otel/sdk` v1.46.0, `otlploggrpc` v0.21.0 and `ebpf` v0.22.0, resolving all of them.

## Impact Summary

The components used by the default `config.yaml` (`otlp` receiver, `batch` processor, `debug` exporter, `health_check` extension) are **not affected** by any config-level breaking change in this range. If you only use those, your existing configuration continues to validate and run unchanged.

The breaking changes below matter only if your configuration uses the affected components. The most important classes of change are:

1. **Silent behavior changes** — `hostmetrics` CPU metrics are no longer reported per logical CPU by default, and the OTTL `set` function no longer ignores `nil` values.
2. **Config-invalidating removals** — the `kafka` exporter/receiver removed several deprecated config fields, and the OTTL `Base64Decode` converter was removed; configs that still use them will be rejected by `validate`.
3. **Internal telemetry changes** — `memory_limiter` metrics were renamed and the histogram buckets of batch size metrics changed. These only affect dashboards/alerts built on the Collector's own metrics.
4. **Removed feature gates** — several stable feature gates were removed. Passing a removed gate id to `--feature-gates` makes the Collector fail at startup.

## Breaking Changes by Component

### Core Collector Changes

#### 1. Go Version Requirement

**v0.160.0 Changes:**
- Minimum Go version increased to 1.26 (core [#15799](https://github.com/open-telemetry/opentelemetry-collector/pull/15799), contrib [#50394](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/50394))
- **Impact**: Only affects users building from source

#### 2. Stable feature gates removed

**v0.155.0 Changes:**
- The following stable feature gates were removed. The gated behavior was already permanent, so nothing changes at runtime.
  - `confighttp.framedSnappy` ([#15420](https://github.com/open-telemetry/opentelemetry-collector/pull/15420))
  - `configoptional.AddEnabledField` ([#15421](https://github.com/open-telemetry/opentelemetry-collector/pull/15421))
  - `confmap.newExpandedValueSanitizer` ([#15418](https://github.com/open-telemetry/opentelemetry-collector/pull/15418))
  - `exporter.PersistRequestContext` ([#15424](https://github.com/open-telemetry/opentelemetry-collector/pull/15424))
  - `otelcol.printInitialConfig` ([#15425](https://github.com/open-telemetry/opentelemetry-collector/pull/15425))
  - `telemetry.UseLocalHostAsDefaultMetricsAddress` ([#15419](https://github.com/open-telemetry/opentelemetry-collector/pull/15419))
  - `pdata.enableRefCounting` ([#15426](https://github.com/open-telemetry/opentelemetry-collector/pull/15426))
  - `processor.tailsamplingprocessor.disableinvertdecisions` ([#48976](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/48976))
  - `filelog.decompressFingerprint` ([#48980](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/48980))
- **Action Required (conditional)**: Remove these ids from any `--feature-gates` flag (command line, systemd unit, container args). An unknown gate id is a startup error.

#### 3. Internal telemetry: batch size histogram buckets changed

**v0.157.0 Changes:**
- Histogram bucket boundaries for `otelcol_exporter_queue_batch_send_size_bytes` and `otelcol_processor_batch_batch_send_size_bytes` were replaced with a power-of-2 byte-scale set spanning 128 B to 16 MiB ([#15535](https://github.com/open-telemetry/opentelemetry-collector/issues/15535))
- **Impact**: Metric names are unchanged; only the `le` buckets differ. Only affects dashboards/alerts on the Collector's **own** internal metrics that hard-code specific `le` values.

#### 4. HTTP client/server keepalive settings moved (deprecation)

**v0.160.0 Changes:**
- Flat keepalive fields of HTTP clients and servers are deprecated in favor of a new `keepalive` section ([#14020](https://github.com/open-telemetry/opentelemetry-collector/issues/14020), [#15308](https://github.com/open-telemetry/opentelemetry-collector/pull/15308))
  - Client: `idle_conn_timeout`, `max_idle_conns`, `max_idle_conns_per_host` → `keepalive::*`; `disable_keep_alives: true` → `keepalive::enabled: false`
  - Server: `idle_timeout` → `keepalive::idle_timeout`; `keep_alives_enabled: false` → `keepalive::enabled: false`
- **Impact**: Applies to every HTTP-based component (`otlp_http`, `prometheus_remote_write`, `elasticsearch`, the HTTP server of the `otlp` receiver, `health_check`, ...). The old fields keep working with a deprecation warning. Setting the old fields **together with** the new `keepalive` section is an error.
- **Action Required**: None for now. Migrate to the `keepalive` section when convenient.

### Components with Breaking Changes

#### 1. hostmetricsreceiver

**v0.157.0 Changes:**
- The `cpu` attribute of `system.cpu.time` and `system.cpu.utilization` is now opt-in ([#49161](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/49161))
- **Impact**: This is a **silent behavior change**. Both metrics are now aggregated across logical CPUs and no longer carry the `cpu` attribute. Per-core series disappear, and queries such as `sum by (cpu)` or `avg by (cpu)` stop working.
- **Action Required**: To restore the previous per-logical-CPU output, enable the attribute explicitly:
  ```yaml
  receivers:
    hostmetrics:
      scrapers:
        cpu:
          metrics:
            system.cpu.time:
              attributes: [cpu, state]
            system.cpu.utilization:
              attributes: [cpu, state]
  ```
- The `system.cpu.logical.count` metric is now enabled by default ([#49325](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/49325)). Set `metrics::system.cpu.logical.count::enabled: false` under the `cpu` scraper if it is not wanted.

#### 2. transformprocessor / filterprocessor (OTTL)

These changes also apply to OTTL-based policies of the `tail_sampling` processor.

**v0.157.0 Changes:**
- The `processor.filter.defaultErrorModeIgnore` and `processor.transform.defaultErrorModeIgnore` feature gates were promoted to stable ([#47232](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/47232), [#47231](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/47231))
- **Impact**: The default `error_mode: ignore` introduced in v0.7.x is now permanent. The gates can no longer be disabled; `--feature-gates=-processor.filter.defaultErrorModeIgnore` / `--feature-gates=-processor.transform.defaultErrorModeIgnore` are rejected at startup.
- **Action Required**: If you opted out via the feature gates, remove them and set `error_mode: propagate` explicitly in the processor config instead.

**v0.161.0 Changes:**
- The `ottl.set.allowNil` feature gate was promoted to beta and is enabled by default ([#50872](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/50872), [#48714](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/48714))
- **Impact**: This is a **silent behavior change**. `set(target, value)` with a `value` that evaluates to `nil` used to be a no-op. It now passes `nil` to the target, which clears it (or errors, depending on the path type). A typical affected statement is `set(attributes["x"], attributes["missing"])`.
- **Action Required**: Review `set` statements whose source may be missing and guard them, for example `set(attributes["x"], attributes["y"]) where attributes["y"] != nil`. To temporarily restore the previous behavior, use `--feature-gates=-ottl.set.allowNil`. The gate is expected to become stable (no opt-out) in a future release.

- The deprecated `Base64Decode` converter was removed ([#50875](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/50875))
- **Impact**: Configs that use `Base64Decode(...)` are **rejected by validation**.
- **Action Required**: Use `Decode(value, "base64")` instead.

#### 3. kafkaexporter / kafkareceiver

Both components share the same client configuration, so all of these changes apply to the exporter and the receiver. Configurations that still set the removed options will be **rejected by validation**.

**v0.160.0 Changes:**
- Removed deprecated `auth::tls` and `auth::plain_text` ([#50202](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/50202))
- **Action Required**: Use the top-level `tls` section and `auth::sasl` with `mechanism: PLAIN` instead.
- **Example**:
  ```yaml
  # Old configuration (v0.7.x)
  exporters:
    kafka:
      auth:
        plain_text:
          username: user
          password: pass
        tls:
          insecure: false

  # New configuration (v0.8.x)
  exporters:
    kafka:
      tls:
        insecure: false
      auth:
        sasl:
          mechanism: PLAIN
          username: user
          password: pass
  ```
- Removed deprecated `resolve_canonical_bootstrap_servers_only`, `auth::sasl::version` and `group_rebalance_strategy` ([#50381](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/50381))
- **Action Required**: Remove the first two fields (they were already no-ops). Replace `group_rebalance_strategy: x` with `group_rebalance_strategies: [x]` in the receiver config.

**v0.161.0 Changes:**
- Configurations that set both `auth::sasl` and `auth::kerberos` are now rejected ([#50748](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/50748))
- **Action Required**: Keep only one of the two.

#### 4. memorylimiterprocessor

**v0.155.0 Changes:**
- Internal metrics were renamed to include the `memory_limiter` prefix ([#11203](https://github.com/open-telemetry/opentelemetry-collector/issues/11203), [#15167](https://github.com/open-telemetry/opentelemetry-collector/pull/15167))
  - `otelcol_processor_accepted_{spans,metric_points,log_records}` → `otelcol_processor_memory_limiter_accepted_{spans,metric_points,log_records}`
  - `otelcol_processor_refused_{spans,metric_points,log_records}` → `otelcol_processor_memory_limiter_refused_{spans,metric_points,log_records}`
- **Impact**: No config change. The old names are no longer emitted. Only affects dashboards/alerts on the Collector's **own** internal metrics (including those collected by the `selfmetrics` receiver).

#### 5. filelogreceiver

**v0.156.0 Changes:**
- The `filelog.protobufCheckpointEncoding` feature gate was promoted to beta ([#49387](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/49387))
- **Impact**: Checkpoints persisted via a storage extension (`file_storage`) are now written in protobuf instead of JSON. Existing JSON checkpoints are still read, so no offsets are lost on upgrade. To opt out, use `--feature-gates=-filelog.protobufCheckpointEncoding`.

**v0.159.0 Changes:**
- `ordering_criteria::top_n: 0` now means "match all files" instead of behaving like `top_n: 1`. Relying on the implicit default `top_n: 1` when `ordering_criteria::sort_by` is set is deprecated ([#47444](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/47444))
- **Action Required (conditional)**: If you use `ordering_criteria`, set `top_n` explicitly.

#### 6. deltatocumulativeprocessor

**v0.158.0 Changes:**
- The component type was renamed from `deltatocumulative` to `delta_to_cumulative` ([#45339](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/45339), [#49634](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/49634))
- **Impact**: The old `deltatocumulative` id remains available as a deprecated alias and logs a warning at startup.
- **Action Required**: Rename `deltatocumulative` to `delta_to_cumulative` in your config when convenient.

> Note: Several other bundled components were already renamed before v0.7.x and log the same kind of warning when the old id is used: `filelog` → `file_log`, `fluentforward` → `fluent_forward`, `hostmetrics` → `host_metrics`, `windowseventlog` → `windows_event_log`, `resourcedetection` → `resource_detection`, `otlp` (exporter) → `otlp_grpc`, `otlphttp` → `otlp_http`, `prometheusremotewrite` → `prometheus_remote_write`. The old ids still work.

#### 7. prometheusreceiver

**v0.156.0 Changes:**
- The `receiver.prometheusreceiver.IgnoreScopeInfoMetric` feature gate was promoted to beta ([#47312](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/47312))
- **Impact**: The `otel_scope_info` metric is no longer used to extract scope attributes. Only matters when scraping endpoints exposed by OpenTelemetry SDKs that emit `otel_scope_info`. To temporarily restore the previous behavior, use `--feature-gates=-receiver.prometheusreceiver.IgnoreScopeInfoMetric`.

#### 8. prometheusremotewriteexporter (deprecation)

**v0.160.0 Changes:**
- `resource_to_telemetry_conversion` is deprecated in favor of `resource_constant_labels` ([#48862](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/48862))
- **Impact**: The old section keeps working with a deprecation warning. It cannot be combined with `resource_constant_labels`.
- **Action Required**: Migrate when convenient.
- **Example**:
  ```yaml
  # Old configuration (v0.7.x)
  exporters:
    prometheusremotewrite:
      resource_to_telemetry_conversion:
        enabled: true

  # New configuration (v0.8.x)
  exporters:
    prometheus_remote_write:
      resource_constant_labels:
        included: ["*"]
  ```

The `sacloud` exporter has already been migrated internally; no configuration change is needed for it.

#### 9. resourcedetectionprocessor (deprecation)

**v0.158.0 Changes:**
- The per-detector `fail_on_missing_metadata` option (e.g. `ec2::fail_on_missing_metadata`) is deprecated in favor of the top-level `fail_on_missing_metadata` of the processor ([#46659](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/46659))
- **Impact**: The old field keeps working with a deprecation warning.

#### 10. tailsamplingprocessor

**v0.157.0 Changes:**
- The unit of the `otelcol_processor_tail_sampling_sampling_policy_execution_time_sum` internal metric changed from `µs` to `us` to comply with UCUM ([#49453](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/49453))
- **Impact**: Internal telemetry metadata only. See also the removed `disableinvertdecisions` feature gate and the OTTL changes above.

### Components with No Breaking Changes

The following bundled components had no component-specific breaking changes in the v0.155.0–v0.161.0 range:
`otlpreceiver`, `otlpexporter`, `otlphttpexporter`, `debugexporter`, `batchprocessor`, `awss3exporter`, `fileexporter`, `elasticsearchexporter`, `healthcheckextension`, `filestorage`, `attributesprocessor`, `groupbyattrsprocessor`, `resourceprocessor`, `jaegerreceiver`, `fluentforwardreceiver`, `journaldreceiver`, `windowseventlogreceiver`.

> Note: `mackerelotlpexporter` is a third-party component (`github.com/mackerelio/opentelemetry-collector-mackerel`) with its own versioning and was not covered by the Collector changelogs above. Review its release notes separately if you use it.

## Configuration Validation

After updating configurations, validate them using:

```bash
./sacloud-otel-collector validate --config your-config.yaml
```

The `kafka` field removals and the `Base64Decode` removal are the changes most likely to make an existing `config.yaml` fail validation outright. The `hostmetrics` CPU attribute change and the OTTL `set` nil handling are the ones most likely to silently change behavior while still validating.

## References

### OpenTelemetry Collector Releases
- [v0.155.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.155.0)
- [v0.156.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.156.0)
- [v0.157.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.157.0)
- [v0.158.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.158.0)
- [v0.159.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.159.0)
- [v0.160.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.160.0)
- [v0.161.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.161.0)

### OpenTelemetry Collector Contrib Releases
- [v0.155.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.155.0)
- [v0.156.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.156.0)
- [v0.157.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.157.0)
- [v0.158.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.158.0)
- [v0.159.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.159.0)
- [v0.160.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.160.0)
- [v0.161.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.161.0)

## Support

If you encounter issues during the upgrade:

1. Check component-specific documentation at [OpenTelemetry Docs](https://opentelemetry.io/docs/)
2. Review the full changelog for each version
3. Report issues at the [sacloud-otel-collector repository](https://github.com/sacloud/sacloud-otel-collector/issues)
