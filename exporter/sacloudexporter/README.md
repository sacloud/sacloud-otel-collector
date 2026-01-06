# SAKURA Cloud Exporter

<!-- status autance added by mdatagen. Do not edit. -->
| Status                   |                                              |
| ------------------------ | -------------------------------------------- |
| Stability                | [development]: metrics, logs, traces         |
| Distributions            | [sacloud-otel-collector]                     |

[development]: https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/component-stability.md#development
[sacloud-otel-collector]: https://github.com/sacloud/sacloud-otel-collector

<!-- end autance added section -->

## Overview

This exporter sends telemetry data to [SAKURA Cloud Monitoring Suite](https://manual.sakura.ad.jp/cloud/appliance/monitoring-suite/).

Internally, it uses the following exporters:
- **Metrics**: [prometheusremotewriteexporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/prometheusremotewriteexporter) (Prometheus Remote Write protocol)
- **Logs/Traces**: [otlphttpexporter](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlphttpexporter) (OTLP/HTTP protocol)

## Configuration

### Required Settings

At least one of the following signal configurations must be specified:

| Field | Description |
|-------|-------------|
| `metrics.endpoint` | Endpoint identifier or full URL for metrics |
| `metrics.token` | Bearer token for metrics authentication |
| `logs.endpoint` | Endpoint identifier or full URL for logs |
| `logs.token` | Bearer token for logs authentication |
| `traces.endpoint` | Endpoint identifier or full URL for traces |
| `traces.token` | Bearer token for traces authentication |

The `endpoint` field accepts either:
- An endpoint identifier from SAKURA Cloud control panel (e.g., `"123456789012"`)
- A full URL (e.g., `"https://123456789012.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write"`)

If only an identifier is provided, it will be expanded to the full URL automatically.

### Optional Settings

| Field | Default | Description |
|-------|---------|-------------|
| `timeout` | `30s` | HTTP request timeout |
| `retry_on_failure.enabled` | `true` | Enable retry on failure |
| `retry_on_failure.initial_interval` | `5s` | Initial retry interval |
| `retry_on_failure.max_interval` | `30s` | Maximum retry interval |
| `retry_on_failure.max_elapsed_time` | `5m` | Maximum elapsed time for retries |

## Example

### Basic configuration with endpoint identifiers

```yaml
exporters:
  sacloud:
    metrics:
      endpoint: "123456789012"
      token: "${SACLOUD_METRICS_TOKEN}"
    logs:
      endpoint: "123456789012"
      token: "${SACLOUD_LOGS_TOKEN}"
    traces:
      endpoint: "123456789012"
      token: "${SACLOUD_TRACES_TOKEN}"
```

### Metrics only

```yaml
exporters:
  sacloud:
    metrics:
      endpoint: "123456789012"
      token: "${SACLOUD_METRICS_TOKEN}"
```

### With custom timeout and retry settings

```yaml
exporters:
  sacloud:
    timeout: 60s
    retry_on_failure:
      enabled: true
      initial_interval: 10s
      max_interval: 60s
      max_elapsed_time: 10m
    metrics:
      endpoint: "123456789012"
      token: "${SACLOUD_METRICS_TOKEN}"
```

## Internal Behavior

### Metrics
- Uses Prometheus Remote Write protocol with Snappy compression
- Resource attributes are converted to metric labels (`resource_to_telemetry_settings.enabled: true`)
- Default queue size: 10,000 items
- Default batch size: 4 MiB per request

### Logs and Traces
- Uses OTLP/HTTP protocol with gzip compression
- Default queue size: 10 MiB buffer
- Default batch size: 4 MiB per request
- Default batch flush timeout: 10 seconds

These defaults are configured to operate safely within SAKURA Cloud Monitoring Suite limits (5 MB per request, 50 requests/second).
