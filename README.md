# sacloud-otel-collector

OpenTelemetry collector for [sacloud](https://github.com/sacloud).

## Usage

### Release Binary

Pre-built binaries are available on [GitHub Releases](https://github.com/sacloud/sacloud-otel-collector/releases):

1. Download the tar.gz archive for your platform from the latest release
2. Extract the archive: `tar xzf sacloud-otel-collector_*.tar.gz`
3. Run with your configuration:

```bash
./sacloud-otel-collector --config config.yml
```

### Container Image

Container images are available on GitHub Container Registry:

```bash
docker pull ghcr.io/sacloud/sacloud-otel-collector:latest
```

Run the collector with your configuration:

```bash
docker run --rm -v $(pwd)/config.yml:/config.yml ghcr.io/sacloud/sacloud-otel-collector:latest --config /config.yml
```

### RPM/DEB Package

RPM and DEB packages are available on [GitHub Releases](https://github.com/sacloud/sacloud-otel-collector/releases).

#### Installation

For Debian/Ubuntu:
```bash
wget https://github.com/sacloud/sacloud-otel-collector/releases/download/v<version>/sacloud-otel-collector_<version>_amd64.deb
sudo dpkg -i sacloud-otel-collector_<version>_amd64.deb
```

For RHEL/CentOS/Fedora:
```bash
wget https://github.com/sacloud/sacloud-otel-collector/releases/download/v<version>/sacloud-otel-collector_<version>_amd64.rpm
sudo rpm -i sacloud-otel-collector_<version>_amd64.rpm
```

#### Configuration

Edit the configuration file at `/etc/sacloud-otel-collector/config.yaml`:

```bash
sudo vi /etc/sacloud-otel-collector/config.yaml
```

The default configuration includes receivers for OTLP, host metrics, and file logs, with exporters configured for SAKURA Cloud Monitoring Suite. You need to replace the `****` placeholders with your actual monitoring suite credentials.

#### Service Management

Enable and start the service:

```bash
sudo systemctl enable sacloud-otel-collector
sudo systemctl start sacloud-otel-collector
```

Check service status:

```bash
sudo systemctl status sacloud-otel-collector
```

View logs:

```bash
sudo journalctl -u sacloud-otel-collector -f
```

Reload configuration after changes:

```bash
sudo systemctl reload sacloud-otel-collector
```

## Components

The sacloud-otel-collector includes the following OpenTelemetry components:

For more details, see [builder-config.yaml](builder-config.yaml).

### Receivers

| Component | Description | Documentation |
|-----------|-------------|---------------|
| otlp | OpenTelemetry Protocol receiver | [Documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/main/receiver/otlpreceiver) |
| prometheus | Prometheus metrics receiver | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/prometheusreceiver) |
| hostmetrics | Host metrics receiver | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/hostmetricsreceiver) |
| jaeger | Jaeger traces receiver | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/jaegerreceiver) |
| kafka | Kafka receiver | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/kafkareceiver) |
| filelog | File log receiver | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/filelogreceiver) |
| fluentforward | Fluent Forward receiver | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/fluentforwardreceiver) |
| journald | Journald receiver | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/journaldreceiver) |

### Processors

| Component | Description | Documentation |
|-----------|-------------|---------------|
| batch | Batch processor | [Documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/batchprocessor) |
| memorylimiter | Memory limiter processor | [Documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/memorylimiterprocessor) |
| resource | Resource processor | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourceprocessor) |
| attributes | Attributes processor | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/attributesprocessor) |
| transform | Transform processor | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/transformprocessor) |
| filter | Filter processor | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/filterprocessor) |
| deltatocumulative | Delta to cumulative processor | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/deltatocumulativeprocessor) |
| resourcedetection | Resource detection processor | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/resourcedetectionprocessor) |
| tailsampling | Tail sampling processor | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/tailsamplingprocessor) |

### Exporters

| Component | Description | Documentation |
|-----------|-------------|---------------|
| debug | Debug exporter | [Documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/debugexporter) |
| otlp | OpenTelemetry Protocol exporter | [Documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlpexporter) |
| prometheusremotewrite | Prometheus Remote Write exporter | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/prometheusremotewriteexporter) |
| otlphttp | OTLP HTTP exporter | [Documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter/otlphttpexporter) |
| file | File exporter | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/fileexporter) |
| elasticsearch | Elasticsearch exporter | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/elasticsearchexporter) |
| awss3 | AWS S3 exporter | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/awss3exporter) |
| sacloud | SAKURA Cloud Monitoring Suite exporter | See [Config examples](#sakuracloud-monitoring-suite) |

### Extensions

| Component | Description | Documentation |
|-----------|-------------|---------------|
| health_check | Health check extension | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/healthcheckextension) |
| file_storage | File storage extension | [Documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/storage/filestorage) |

## Build

See [Building a custom collector](https://opentelemetry.io/docs/collector/custom-collector/).

1. Install Go.
2. `make`
    See [Makefile](Makefile) for details.

`sacloud-otel-collector` will be created.

## Config examples

### Sakuracloud Monitoring Suite

See [manual](https://manual.sakura.ad.jp/cloud/appliance/monitoring-suite/index.html).

#### Using the sacloud exporter (recommended)

The `sacloud` exporter simplifies configuration for SAKURA Cloud Monitoring Suite with sensible defaults.

**Features:**

| Feature | Default | Description |
|---------|---------|-------------|
| Endpoint | - | Endpoint ID (e.g., `123456789012`) or full URL |
| Compression | Enabled | snappy for metrics, gzip for logs/traces |
| Timeout | 30s | HTTP request timeout |
| Retry | Enabled | Exponential backoff (initial: 5s, max: 30s, max elapsed: 5m) |
| Queue | Enabled | Optimized for 5MB per request limit |

**Queue defaults:**
- **Logs/Traces**: 10MiB buffer, 4MiB max batch size, 2 consumers
- **Metrics**: 10000 queue size, 2 consumers

See also [SAKURA Cloud Monitoring Suite limits](https://manual.sakura.ad.jp/cloud/appliance/monitoring-suite/about.html#monitoring-suite-specification-limit).

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
  hostmetrics:
    collection_interval: 10s
    scrapers:
      cpu:
        metrics:
          system.cpu.utilization:
            enabled: true
      memory:
        metrics:
          system.memory.utilization:
            enabled: true
      disk:
      filesystem:
        metrics:
          system.filesystem.utilization:
            enabled: true
      network:
      paging:
        metrics:
          system.paging.utilization:
            enabled: true
  filelog:
    start_at: end
    exclude: []
    include:
      - /var/log/example.log

processors:
  resourcedetection:
    detectors: [system]
    system:
      hostname_sources: [os]

# Replace endpoint identifiers and tokens with your monitoring suite's configurations.
exporters:
  sacloud:
    metrics:
      endpoint: "123456789012" # or "https://123456789012.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write"
      token: "${SACLOUD_METRICS_TOKEN}" # met-***************
    logs:
      endpoint: "123456789012" # or "https://123456789012.logs.monitoring.global.api.sacloud.jp"
      token: "${SACLOUD_LOGS_TOKEN}" # log-***************
    traces:
      endpoint: "123456789012" # or "https://123456789012.traces.monitoring.global.api.sacloud.jp"
      token: "${SACLOUD_TRACES_TOKEN}" # trc-***************

service:
  pipelines:
    metrics:
      receivers: [hostmetrics]
      processors: [resourcedetection]
      exporters: [sacloud]
    logs:
      receivers: [filelog]
      processors: [resourcedetection]
      exporters: [sacloud]
    traces:
      receivers: [otlp]
      processors: [resourcedetection]
      exporters: [sacloud]
```

##### Advanced: Queue Configuration

The default queue settings work well for most use cases. If you need to customize queue behavior, you can override the defaults:

```yaml
exporters:
  sacloud:
    metrics:
      endpoint: "123456789012"
      token: "${SACLOUD_METRICS_TOKEN}"
      # Override default remote_write_queue for metrics
      remote_write_queue:
        enabled: true
        queue_size: 10000
        num_consumers: 2
    logs:
      endpoint: "123456789012"
      token: "${SACLOUD_LOGS_TOKEN}"
      # Override default sending_queue for logs
      sending_queue:
        enabled: true
        sizer: bytes
        queue_size: 10485760  # 10MiB
        num_consumers: 2
        batch:
          flush_timeout: 10s
          max_size: 4194304   # 4MiB per request
    traces:
      endpoint: "123456789012"
      token: "${SACLOUD_TRACES_TOKEN}"
      # Same options as logs
      sending_queue:
        enabled: true
        sizer: bytes
        queue_size: 10485760
        num_consumers: 2
        batch:
          flush_timeout: 10s
          max_size: 4194304
```

##### Advanced: Timeout and Retry Configuration

The exporter includes sensible defaults for timeout (30 seconds) and retry (enabled with exponential backoff). You can customize these settings if needed:

```yaml
exporters:
  sacloud:
    # Timeout for HTTP requests (default: 30s)
    timeout: 30s
    # Retry configuration (default: enabled with exponential backoff)
    retry_on_failure:
      enabled: true
      initial_interval: 5s    # Time to wait after first failure
      max_interval: 30s       # Maximum backoff interval
      max_elapsed_time: 5m    # Maximum total retry time before giving up
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

#### Using standard exporters

Alternatively, you can use standard exporters directly:

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
  hostmetrics:
    collection_interval: 10s
    scrapers:
      cpu:
        metrics:
          system.cpu.utilization:
            enabled: true
      memory:
        metrics:
          system.memory.utilization:
            enabled: true
      disk:
      filesystem:
        metrics:
          system.filesystem.utilization:
            enabled: true
      network:
      paging:
        metrics:
          system.paging.utilization:
            enabled: true
  filelog:
    start_at: end
    exclude: []
    include:
      - /var/log/example.log

processors:
  resourcedetection:
    detectors: [system]
    system:
      hostname_sources: [os]
  batch:
    timeout: 1s
    # Adjust the send_batch_size/send_batch_max_size so that
    # the payload size of a single request does not exceed 5 MiB.
    send_batch_size: 4096
    send_batch_max_size: 4096

# You should replace `****` to your monitoring suite's configurations.
exporters:
  otlphttp/sakura-monitoring-suite-log:
    endpoint: https://****.logs.monitoring.global.api.sacloud.jp
    headers:
      Authorization: "Bearer ****"
  prometheusremotewrite/sakura-monitoring-suite-metrics:
    endpoint: https://****.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write
    headers:
      Authorization: "Bearer ****"
    resource_to_telemetry_conversion:
      enabled: true
  otlphttp/sakura-monitoring-suite-trace:
    endpoint: https://****.traces.monitoring.global.api.sacloud.jp
    headers:
      Authorization: "Bearer ****"

service:
  pipelines:
    metrics:
      receivers:
        - hostmetrics
      processors:
        - resourcedetection
        - batch
      exporters:
        - prometheusremotewrite/sakura-monitoring-suite-metrics
    logs:
      receivers:
        - filelog
      processors:
        - resourcedetection
        - batch
      exporters:
        - otlphttp/sakura-monitoring-suite-log
    traces:
      receivers:
        - otlp
      processors:
        - resourcedetection
        - batch
      exporters:
        - otlphttp/sakura-monitoring-suite-trace
```

## Contributing

If you want to add a new component of otel collector, modify [builder-config.yaml](builder-config.yaml). Then run `make build-src && make` to build the collector.

## Automation

GitHub Actions are used to automate the build and release process. The following steps are performed:

1. When new branches are merged into the `main` branch, [tagpr](https://github.com/Songmu/tagpr) creates/updates a new pull request for the release.
2. When the release pull request is merged, the following steps are performed:
   1. The `main` branch is tagged with a new version number.
   2. The `main` branch is built and released to GitHub Releases.

## License

sacloud-otel-collector Copyright (C) 2025 The sacloud/sacloud-otel-collector authors.
This project is published under [Apache 2.0 License](LICENSE).
