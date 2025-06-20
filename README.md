# sacloud-otel-collector

OpenTelemetry collector for [sacloud](https://github.com/sacloud).

## Build

See [Building a custom collector](https://opentelemetry.io/docs/collector/custom-collector/).

1. Install Go.
2. `make`
    See [Makefile](Makefile) for details.

`build/sacloud-otel-collector` will be created.

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
