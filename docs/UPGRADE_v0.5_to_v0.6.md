# Upgrade Guide: v0.5 to v0.6

This document outlines the upgrade considerations when migrating from sacloud-otel-collector v0.5.x to v0.6.x.

## OpenTelemetry Collector Version Update

- **v0.5.x**: Based on OpenTelemetry Collector v0.142.0
- **v0.6.x**: Based on OpenTelemetry Collector v0.142.0 (unchanged)

The bundled components and their versions are the same as in v0.5.x. No configuration changes are required.

## Changes

### Build Toolchain

#### 1. Release binaries are built with Go 1.26

Release binaries, packages and container images are now built with Go 1.26 instead of Go 1.25 ([#83](https://github.com/sacloud/sacloud-otel-collector/pull/83)).

- **Impact**: No configuration or behavior change is expected. Building from source is still tested with both Go 1.25 and Go 1.26
- **Action Required**: None

### Documentation

#### 1. Package download URLs in README corrected

The RPM/DEB download URLs in the README were missing the `_linux` part of the file name ([#79](https://github.com/sacloud/sacloud-otel-collector/pull/79)). The package file names themselves did not change:

```
sacloud-otel-collector_<version>_linux_amd64.deb
sacloud-otel-collector_<version>_linux_amd64.rpm
```

- **Action Required**: If you scripted the download based on the old README, use the file names above

## Configuration Validation

Validate your configuration using:

```bash
./sacloud-otel-collector validate --config your-config.yaml
```

## Support

If you encounter issues during the upgrade, report them at the [sacloud-otel-collector repository](https://github.com/sacloud/sacloud-otel-collector/issues).
