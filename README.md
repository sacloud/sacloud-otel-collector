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
