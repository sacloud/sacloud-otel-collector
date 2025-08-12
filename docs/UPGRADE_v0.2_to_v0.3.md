# Upgrade Guide: v0.2 to v0.3

This document outlines the breaking changes and upgrade considerations when migrating from sacloud-otel-collector v0.2.x (based on OpenTelemetry Collector v0.125.0) to v0.3.x (based on OpenTelemetry Collector v0.131.0).

## OpenTelemetry Collector Version Update

- **v0.2.x**: Based on OpenTelemetry Collector v0.125.0
- **v0.3.x**: Based on OpenTelemetry Collector v0.131.0

## Breaking Changes by Component

### Components with Breaking Changes

#### 1. **otlpexporter / otlphttpexporter**

**v0.130.0 Changes:**
- `exporter/otlp`: Remove deprecated batcher config from OTLP, use queuebatch ([#13339](https://github.com/open-telemetry/opentelemetry-collector/pull/13339))
- **Action Required**: If using `batcher` configuration, migrate to `queuebatch`
- **Example**:
  ```yaml
  # Old configuration (v0.2.x)
  exporters:
    otlp:
      batcher:
        enabled: true
  
  # New configuration (v0.3.x)
  exporters:
    otlp:
      queuebatch:
        enabled: true
  ```

**v0.131.0 Changes:**
- `confighttp`: Move `confighttp.framedSnappy` feature gate to beta ([#10584](https://github.com/open-telemetry/opentelemetry-collector/issues/10584))
- **Impact**: May affect HTTP/gRPC communication with snappy compression

#### 2. **otlpreceiver**

**v0.128.0 Changes:**

- `otlpreceiver`: Use configoptional.Optional to define optional configuration sections in the OTLP receiver. Remove Unmarshal method. ([#13119](https://github.com/open-telemetry/opentelemetry-collector/pull/13119))

#### 3. **prometheusreceiver**

**v0.126.0 Changes:**
- `receiver/prometheus`: Upgrade the RemoveLegacyResourceAttributes featuregate to beta ([#32814](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/32814))

**v0.129.0 Changes:**
- `receiver/prometheus`: Promote the receiver.prometheusreceiver.RemoveLegacyResourceAttributes featuregate to stable ([#40572](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/40572))
- **Impact**: Legacy attributes automatically converted to new semantic conventions
  - `net.host.name` → `server.address`
  - `net.host.port` → `server.port`
  - `http.scheme` → `url.scheme`
- **Note**: This is automatic and requires no configuration changes, but may affect downstream processing

#### 4. **prometheusremotewriteexporter**

**v0.129.0 Changes:**
- `exporter/prometheusremotewrite`: Remove the stable exporter.prometheusremotewriteexporter.deprecateCreatedMetric featuregate ([#40570](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/40570))

#### 5. **kafkareceiver**

**v0.129.0 Changes:**
- `receiver/kafka`: Improve kafkareceiver internal metrics telemetry ([#40816](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/40816))
  - Added new metrics: kafka_broker_connects, kafka_broker_closed
  - Removed explicit component "name" attribute
  - Changed "partition" attribute type to int64

**v0.130.0 Changes:**
- `kafka`: Default client ID now honours configuration and defaults to `otel-collector` (previously always `sarama`) ([#41090](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/41090))
- **Action Required**: If relying on specific client ID for monitoring or ACLs, update configurations

#### 6. **elasticsearchexporter**

**v0.129.0 Changes:**
- `exporter/elasticsearch`: Add better ECS mapping for traces when using ECS mapping mode ([#40807](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/40807))
- **Impact**: If using ECS mapping mode, trace processing behavior will change (improvement)

#### 7. **deltatocumulativeprocessor**

**v0.131.0 Changes:**
- `deltatocumulativeprocessor`: Unexport Processor, CountingSink ([#40656](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/40656))
- **Impact**: Only affects custom code using the processor's API, not configuration

### Components Without Breaking Changes

The following components have no breaking changes between v0.125.0 and v0.131.0:

- **Exporters**: fileexporter
- **Receivers**: filelogreceiver, jaegerreceiver, hostmetricsreceiver
- **Processors**: batchprocessor, memorylimiterprocessor, resourceprocessor, attributesprocessor, transformprocessor, filterprocessor
- **Extensions**: healthcheckextension, filestorage

## Configuration Validation

After updating configurations, validate them using:

```bash
./sacloud-otel-collector validate --config your-config.yaml
```

## References

### OpenTelemetry Collector Releases
- [v0.126.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.126.0)
- [v0.127.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.127.0)
- [v0.128.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.128.0)
- [v0.129.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.129.0)
- [v0.130.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.130.0)
- [v0.131.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.131.0)

### OpenTelemetry Collector Contrib Releases
- [v0.126.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.126.0)
- [v0.127.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.127.0)
- [v0.128.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.128.0)
- [v0.129.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.129.0)
- [v0.130.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.130.0)
- [v0.131.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.131.0)

## Support

If you encounter issues during the upgrade:

1. Check component-specific documentation at [OpenTelemetry Docs](https://opentelemetry.io/docs/)
2. Review the full changelog for each version
3. Report issues at the [sacloud-otel-collector repository](https://github.com/sacloud/sacloud-otel-collector/issues)
