# Upgrade Guide: v0.3 to v0.4

This document outlines the upgrade considerations when migrating from sacloud-otel-collector v0.3.x to v0.4.x.

## OpenTelemetry Collector Version Update

- **v0.3.x**: Based on OpenTelemetry Collector v0.131.0
- **v0.4.x**: Based on OpenTelemetry Collector v0.131.0 (unchanged)

The bundled components and their versions are the same as in v0.3.x. No configuration changes are required.

## Changes

### Container Image

#### 1. Base image changed from `scratch` to distroless

The base image of `ghcr.io/sacloud/sacloud-otel-collector` changed from `scratch` to `gcr.io/distroless/static-debian12` ([#60](https://github.com/sacloud/sacloud-otel-collector/pull/60)).

- The image now contains the time zone database (`tzdata`), which the `scratch` based image lacked, so time zones can now be resolved by name inside the container
- CA certificates are provided by the base image instead of being copied from Alpine
- The binary path (`/sacloud-otel-collector`), the default configuration path (`/etc/otelcol-contrib/config.yaml`) and the user (`10001:10001`) are unchanged
- **Impact**: Only affects users of the container image. The image still has no shell or package manager
- **Action Required**: None in most cases. If you build your own image on top of this one and rely on the root filesystem containing only the binary, the configuration and the CA certificates, review your Dockerfile

### Binary and Packages

No changes. RPM/DEB packages and release binaries behave the same as in v0.3.x.

## Configuration Validation

Validate your configuration using:

```bash
./sacloud-otel-collector validate --config your-config.yaml
```

## Support

If you encounter issues during the upgrade, report them at the [sacloud-otel-collector repository](https://github.com/sacloud/sacloud-otel-collector/issues).
