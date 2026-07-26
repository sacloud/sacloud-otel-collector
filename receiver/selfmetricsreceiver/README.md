# Selfmetrics Receiver

The `selfmetrics` receiver collects the collector's own internal telemetry
metrics (e.g. `otelcol_exporter_sent_metric_points`, `otelcol_exporter_queue_size`,
`otelcol_process_memory_rss`) by scraping the collector's internal telemetry
Prometheus endpoint, so they can be sent through a normal metrics pipeline to
any configured exporter.

It is a thin wrapper around the
[Prometheus receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/prometheusreceiver)
with a scrape configuration targeting the collector's own
`service::telemetry::metrics` endpoint (default `127.0.0.1:8888`) baked in.

## Configuration

No configuration is required. The receiver works out of the box with the
collector's default internal telemetry settings:

```yaml
receivers:
  selfmetrics:

service:
  pipelines:
    metrics:
      receivers: [selfmetrics]
      exporters: [sacloud]
```

### Optional settings

| Name | Default | Description |
|------|---------|-------------|
| `endpoint` | `127.0.0.1:8888` | `host:port` of the internal telemetry Prometheus endpoint. Change this only if `service::telemetry::metrics` is configured to use a different address. |
| `collection_interval` | `60s` | Scrape interval. |

Example with custom settings:

```yaml
receivers:
  selfmetrics:
    endpoint: 127.0.0.1:18888
    collection_interval: 30s

service:
  telemetry:
    metrics:
      readers:
        - pull:
            exporter:
              prometheus:
                host: "127.0.0.1"
                port: 18888
```

## Notes

- The `job` label (and thus the `service.name` resource attribute) of the
  collected metrics is `sacloud-otel-collector`.
- The `service.instance.id` resource attribute is derived from the scrape
  target address, which is the same on every host by default. When running
  multiple collectors, add a processor such as `resourcedetection` to the
  pipeline to distinguish hosts.
- To get more detailed metrics (e.g. batch processor internals), raise the
  telemetry level:

  ```yaml
  service:
    telemetry:
      metrics:
        level: detailed
  ```
