# Upgrade Guide: v0.4 to v0.5

This document outlines the breaking changes and upgrade considerations when migrating from sacloud-otel-collector v0.4.x (based on OpenTelemetry Collector v0.131.0) to v0.5.x (based on OpenTelemetry Collector v0.142.0).

## OpenTelemetry Collector Version Update

- **v0.4.x**: Based on OpenTelemetry Collector v0.131.0
- **v0.5.x**: Based on OpenTelemetry Collector v0.142.0

## Breaking Changes by Component

### Core Collector Changes

#### 1. Go Version Requirement

**v0.133.0 Changes:**
- Minimum Go version increased to 1.24
- **Impact**: Only affects users building from source

### Components with Breaking Changes

#### 1. prometheusreceiver

**v0.142.0 Changes:**
- `receiver.prometheusreceiver.RemoveStartTimeAdjustment` feature gate promoted to stable ([#44180](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/44180))
- In-receiver metric start time adjustment removed
- **Action Required**: Migrate to the `metricstarttime` processor if you relied on start time adjustment
- **Deprecated**: `use_start_time_metric` and `start_time_metric_regex` config options
- Native histogram scraping and ingestion is controlled by `scrape_native_histograms` in the Prometheus scrape configuration
  - The feature gate `receiver.prometheusreceiver.EnableNativeHistograms` is now stable and enabled by default
  - **Action Required**: To actually scrape native histograms, set `scrape_native_histograms: true` and include `PrometheusProto` in `scrape_protocols`

#### 2. elasticsearchexporter

**v0.132.0 Changes:**
- Default `flush::interval` and `batcher::flush_timeout` changed to 10s ([#41726](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/41726))
- `batcher` config deprecated in favor of `sending_queue`

**v0.138.0 Changes:**
- `batcher` configuration removed entirely ([#42767](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/42767))
- `num_consumers` and `flush` deprecated as they conflict with `sending_queue`
- **Action Required**: Migrate from `batcher` to `sending_queue` configuration
- **Example**:
  ```yaml
  # Old configuration (v0.4.x)
  exporters:
    elasticsearch:
      batcher:
        enabled: true
        flush_timeout: 5s

  # New configuration (v0.5.x)
  exporters:
    elasticsearch:
      sending_queue:
        enabled: true
  ```

#### 3. kafkareceiver

**v0.137.0 Changes:**
- Franz-go is now the default Kafka client library
- Added max_partition_fetch_size option

**v0.141.0 Changes:**
- Deprecated `topic` and `encoding` options removed ([#44568](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/44568))
- `receiver.kafkareceiver.UseFranzGo` feature gate moved to Stable ([#44598](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/44598))
- Franz-go is now the default and only Kafka client
- Sarama client will be removed after v0.143.0

**v0.142.0 Changes:**

- `default_fetch_size` deprecated; use `max_fetch_size` instead
- Singular `topic` and `exclude_topic` options deprecated in favor of list-based configuration (`topics`, `exclude_topics`)
- Example:
```yaml
# Old configuration (v0.4.x)
receivers:
  kafka:
    logs:
      topic: my-topic
      exclude_topic: excluded-topic

# New configuration (v0.5.x)
receivers:
  kafka:
    logs:
      topics:
        - my-topic
      exclude_topics:
        - excluded-topic
```

#### 4. resourcedetectionprocessor

**v0.142.0 Changes:**

- `processor.resourcedetection.propagateerrors` feature gate promoted to beta (now enabled by default)
  - **Impact**: Resource detector errors now stop the collector from starting (previously errors were only logged)
  - To restore previous behavior: `--feature-gates=-processor.resourcedetection.propagateerrors`
- Deprecated `attributes` configuration option removed ([#44616](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/44616))
- **Action Required**: If using `attributes` config, migrate to the standard resource attribute configuration

## Configuration Validation

After updating configurations, validate them using:

```bash
./sacloud-otel-collector validate --config your-config.yaml
```

## References

### OpenTelemetry Collector Releases
- [v0.132.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.132.0)
- [v0.133.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.133.0)
- [v0.134.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.134.0)
- [v0.135.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.135.0)
- [v0.136.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.136.0)
- [v0.137.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.137.0)
- [v0.138.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.138.0)
- [v0.139.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.139.0)
- [v0.140.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.140.0)
- [v0.141.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.141.0)
- [v0.142.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.142.0)

### OpenTelemetry Collector Contrib Releases
- [v0.132.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.132.0)
- [v0.133.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.133.0)
- [v0.134.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.134.0)
- [v0.135.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.135.0)
- [v0.136.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.136.0)
- [v0.137.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.137.0)
- [v0.138.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.138.0)
- [v0.139.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.139.0)
- [v0.140.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.140.0)
- [v0.141.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.141.0)
- [v0.142.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.142.0)

## Support

If you encounter issues during the upgrade:

1. Check component-specific documentation at [OpenTelemetry Docs](https://opentelemetry.io/docs/)
2. Review the full changelog for each version
3. Report issues at the [sacloud-otel-collector repository](https://github.com/sacloud/sacloud-otel-collector/issues)
