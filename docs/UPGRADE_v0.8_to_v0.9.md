# Upgrade Guide: v0.8 to v0.9

This document outlines the breaking changes and upgrade considerations when migrating from sacloud-otel-collector v0.8.x (based on OpenTelemetry Collector v0.161.0) to v0.9.x (based on OpenTelemetry Collector v0.166.0).

## OpenTelemetry Collector Version Update

- **v0.8.x**: Based on OpenTelemetry Collector v0.161.0
- **v0.9.x**: Based on OpenTelemetry Collector v0.166.0

## Why This Upgrade

This upgrade is driven by a security fix. The OpenTelemetry Go SDK pinned by Collector v0.161.0 has an open security advisory that cannot be patched without moving the Collector forward:

- `go.opentelemetry.io/otel/sdk`, `go.opentelemetry.io/otel/sdk/log` and `go.opentelemetry.io/otel/trace` — `AttributeValueLengthLimit` is bypassed by attribute values composed of U+FFFD runes ([GHSA-q7gw-24fr-5p69](https://github.com/advisories/GHSA-q7gw-24fr-5p69))

v0.166.0 ships these modules at v1.47.0, which resolves it.

## Impact Summary

Whether you are affected depends on which components your configuration uses. Review the section of each component you use under [Breaking Changes by Component](#breaking-changes-by-component); the [Core Collector Changes](#core-collector-changes) apply to every configuration. Components that need no action are listed under [Components with No Breaking Changes](#components-with-no-breaking-changes).

The most important classes of change are:

1. **Config-invalidating removals** — the flat keepalive fields of HTTP clients and servers, deprecated in v0.8.x, were removed. The `elasticsearch` exporter, the `kafka` receiver and the `transform` processor also removed deprecated options. Experimental OTTL converters now require a feature gate, and OTTL rejects non-list literals for list arguments. Configs that still use them will be rejected by `validate`.
2. **Silent behavior changes** — metric label names starting with `_` are no longer prefixed with `key_` (this also applies to the `sacloud` exporter), the `jaeger` receiver emits `http.response.status_code` instead of `http.status_code`, and the `elasticsearch` exporter retries on more errors by default.
3. **New startup warning** — a pipeline that combines the `batch` processor with an exporter that batches by itself (including the `sacloud` exporter for logs and traces) now logs a double batching warning.
4. **Removed feature gates** — several stable feature gates were removed. Passing a removed gate id to `--feature-gates` makes the Collector fail at startup.

## Breaking Changes by Component

### Core Collector Changes

#### 1. Flat HTTP client/server keepalive settings removed

**v0.165.0 Changes:**
- The flat keepalive fields deprecated in v0.8.x were removed ([#15973](https://github.com/open-telemetry/opentelemetry-collector/pull/15973))
  - Client: `idle_conn_timeout`, `max_idle_conns`, `max_idle_conns_per_host` → `keepalive::*`; `disable_keep_alives: true` → `keepalive::enabled: false`
  - Server: `idle_timeout` → `keepalive::idle_timeout`; `keep_alives_enabled: false` → `keepalive::enabled: false`
- **Impact**: Applies to every HTTP-based component (`otlp_http`, `prometheus_remote_write`, `elasticsearch`, the HTTP server of the `otlp` receiver, `health_check`, ...). Configs that still set the old fields are **rejected by validation**.
- **Action Required**: Move the old fields into the `keepalive` section.
- **Example**:
  ```yaml
  # Old configuration (v0.8.x)
  exporters:
    otlp_http:
      endpoint: https://example.com
      max_idle_conns: 10
      idle_conn_timeout: 30s

  # New configuration (v0.9.x)
  exporters:
    otlp_http:
      endpoint: https://example.com
      keepalive:
        max_idle_conns: 10
        idle_conn_timeout: 30s
  ```

#### 2. Double batching warning

**v0.165.0 Changes:**
- The Collector now warns at startup when a pipeline has both the `batch` processor and an exporter with batching enabled in its `sending_queue` ([#15965](https://github.com/open-telemetry/opentelemetry-collector/issues/15965)):
  ```
  Pipeline has both a batch processor and an exporter with batching enabled; this will cause double batching
  ```
- **Impact**: Warning only; the pipeline keeps working. The `sacloud` exporter always batches logs and traces by itself, so a logs or traces pipeline that uses it together with the `batch` processor logs this warning. Metrics pipelines are not affected.
- **Action Required**: Remove the `batch` processor from logs and traces pipelines that export with the `sacloud` exporter. It is not needed there.

#### 3. Stable feature gates removed

- The following stable feature gates were removed. The gated behavior was already permanent, so nothing changes at runtime.
  - `ottl.PanicDuplicateName` and `ottl.contexts.enableOTelColContext` (v0.162.0, [#46437](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/46437))
  - `processor.transform.defaultErrorModeIgnore` (v0.162.0, [#51242](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/51242))
  - `processor.filter.defaultErrorModeIgnore` (v0.164.0, [#52238](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/52238))
  - `ottl.set.allowNil` (v0.164.0, [#48714](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/48714))
  - `receiver.prometheusreceiver.IgnoreScopeInfoMetric` (v0.164.0, [#47312](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/47312))
  - `filelog.protobufCheckpointEncoding` (v0.165.0, [#52387](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/52387))
- The following feature gates were promoted to stable and can no longer be disabled.
  - `pkg.translator.prometheus.PermissiveLabelSanitization` (v0.164.0, see [prometheusremotewriteexporter](#1-prometheusremotewriteexporter--sacloudexporter-metrics))
  - `connector.routing.defaultErrorModeIgnore` (v0.165.0, [#52352](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/52352))
- **Action Required (conditional)**: Remove these ids from any `--feature-gates` flag (command line, systemd unit, container args). An unknown gate id, or disabling a stable gate, is a startup error.
- **Impact of the removed opt-outs**: The behaviors that v0.8.x allowed to opt out of are now permanent: OTTL `set` passes `nil` values to the target, the `prometheus` receiver no longer extracts scope attributes from `otel_scope_info`, and `file_log` checkpoints are always written as protobuf (existing JSON checkpoints are still read).

### Components with Breaking Changes

#### 1. prometheusremotewriteexporter / sacloudexporter (metrics)

The `sacloud` exporter sends metrics with the `prometheus_remote_write` exporter, so these changes apply to it as well.

**v0.162.0 / v0.164.0 Changes:**
- Leading and consecutive underscores in label names are always preserved ([#50429](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/50429))
- **Impact**: This is a **silent behavior change**. A label starting with a single underscore (e.g. an attribute named `_foo`) used to be exported as `key_foo`. It is now exported as `_foo`. Queries, dashboards and alerts that refer to the `key_` prefixed label names stop matching, and new time series are created under the new label names. The behavior became the default in v0.162.0 and permanent in v0.164.0, so it cannot be reverted with `--feature-gates=-pkg.translator.prometheus.PermissiveLabelSanitization`.
- **Action Required (conditional)**: If you have metric or resource attributes whose names start with `_`, update the queries that use them.

**v0.164.0 Changes:**
- `timeout` now applies to each attempt to send a request instead of to the whole retry loop ([#51983](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/51983))
- **Impact**: With the standard `prometheus_remote_write` exporter, failed requests are now retried until `retry_on_failure::max_elapsed_time` instead of being dropped after `timeout` (5s by default), and the retry delay grows up to `retry_on_failure::max_interval`. When only one of `timeout` and `http::timeout` is set, the other one takes the same value.
- The `sacloud` exporter already applied its `timeout` to each request in v0.8.2, so no configuration change is needed for it. Its retry delay for metrics now grows up to `retry_on_failure::max_interval` as well.

#### 2. transformprocessor / filterprocessor (OTTL)

These changes also apply to OTTL-based conditions of the `routing` connector and policies of the `tail_sampling` processor.

**v0.162.0 / v0.163.0 Changes:**
- The experimental converters `All`, `Any`, `Filter`, `Find`, `MapEach`, `MapKeys`, `Reduce` and `When` are no longer part of the standard functions ([#51563](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/51563)), and experimental functions must be enabled with a feature gate ([#51564](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/51564))
- **Impact**: Configs that use these converters are **rejected by validation** unless the feature gates are enabled.
- **Action Required (conditional)**: Avoid these converters, or enable them with `--feature-gates=ottl.functions.enableLambda`. Other experimental functions require `pkg.ottl.functions.enableExperimental`. Experimental functions are not covered by OTTL's stability guarantees.

**v0.163.0 Changes:**
- Literal values that are not lists are rejected for list arguments at parse time ([#51833](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/51833))
- **Impact**: Statements such as `keep_keys(attributes, "http.method")` or `Concat("a", "-")` used to pass validation and fail on every execution. They are now **rejected by validation**.
- **Action Required**: Pass a list, for example `keep_keys(attributes, ["http.method"])`.

- `Time` no longer treats a trailing `%Z` as UTC when no location is given ([#51416](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/51416))
- **Impact**: This is a **silent behavior change** for `Time(value, format)` calls whose format ends with `%Z`: the time zone abbreviation in the value is now parsed instead of being assumed to be UTC.
- **Action Required (conditional)**: Review such `Time` calls and pass the location explicitly if a fixed time zone is intended.

**v0.164.0 / v0.165.0 Changes (transform processor):**
- The `flatten_data` option was deprecated in favor of `flatten::log` ([#51766](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/51766)) and then removed ([#32080](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/32080))
- **Impact**: Configs that still set `flatten_data` are **rejected by validation**.
- **Action Required**: Replace `flatten_data: true` with `flatten: {log: true}`. It still requires the `transform.flatten.logs` feature gate.

#### 3. elasticsearchexporter

Configurations that still set the removed options will be **rejected by validation**.

**v0.162.0 Changes:**
- Removed deprecated `flush` and `num_workers` ([#42718](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/42718))
- **Action Required**: Migrate `num_workers` to `sending_queue::num_consumers`, `flush::interval` to `sending_queue::batch::flush_timeout`, and `flush::bytes` to `sending_queue::batch::max_size` with `sending_queue::batch::sizer: bytes`.

**v0.165.0 Changes:**
- Removed deprecated `mapping::mode`, `retry::max_requests`, `logs_dynamic_index`, `metrics_dynamic_index` and `traces_dynamic_index` ([#52027](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/52027))
- **Action Required**: Use `retry::max_retries` instead of `retry::max_requests`, and the `X-Elastic-Mapping-Mode` client metadata key or the `elastic.mapping.mode` scope attribute instead of `mapping::mode`.

**v0.166.0 Changes:**
- The exporter now retries by default on transport errors, client timeouts and bulk HTTP statuses 400, 401, 403, 404, 409, 429, 500, 502, 503 and 504, and on document statuses 429, 500, 502, 503 and 504 ([#52551](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/52551))
- **Impact**: This is a **silent behavior change** towards at-least-once delivery. `retry::retry_on_document_status` no longer follows `retry::retry_on_status`.
- **Action Required (conditional)**: To keep the previous default, set both `retry::retry_on_status` and `retry::retry_on_document_status` to `[429]`. If you set `retry::retry_on_status` only, set `retry::retry_on_document_status` to the same list.

#### 4. kafkareceiver / kafkaexporter

**v0.162.0 Changes:**
- `metadata::retry::max` and `metadata::retry::backoff` are applied again, and the default of `metadata::retry::max` changed from 3 to 20 ([#50116](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/50116))
- **Impact**: These fields were silently ignored before, so the defaults behave as in v0.8.x. Only configurations that set them explicitly change behavior: the configured values are now honored.
- `metadata::full` is deprecated and has no effect ([#50116](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/50116))

**v0.164.0 / v0.166.0 Changes (receiver):**
- `partition_processing::max_buffered_batches` was deprecated ([#52253](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/52253)) and then removed ([#52922](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/52922))
- **Impact**: Configs that still set it are **rejected by validation**.
- **Action Required**: Use `partition_processing::max_prefetch::batches` instead.
- The built-in `jaeger_proto`, `jaeger_json` and `json` (logs) encodings are deprecated in favor of encoding extensions ([#39113](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/39113))
- **Impact**: The encodings keep working. The encoding extensions that replace them are not bundled in sacloud-otel-collector, so no action is possible for now.

#### 5. jaegerreceiver

**v0.165.0 Changes:**
- The `pkg.translator.jaeger.EmitV1HttpConventions` and `pkg.translator.jaeger.DontEmitV0HttpConventions` feature gates were promoted to beta ([#52146](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/52146))
- **Impact**: This is a **silent behavior change**. Spans received in Jaeger format now carry `http.response.status_code` instead of the deprecated `http.status_code` attribute. Processors, sampling policies and queries that refer to `http.status_code` stop matching.
- **Action Required (conditional)**: Update references to `http.status_code`. To temporarily restore the previous behavior, use `--feature-gates=-pkg.translator.jaeger.EmitV1HttpConventions,-pkg.translator.jaeger.DontEmitV0HttpConventions`.

#### 6. prometheusreceiver

**v0.165.0 Changes:**
- `__scheme__` and `__metrics_path__` labels on scraped samples are no longer dropped from datapoint attributes ([#52736](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/52736))
- **Impact**: Only matters when a scraped endpoint exposes series that carry these labels themselves.

See also the removed `receiver.prometheusreceiver.IgnoreScopeInfoMetric` feature gate above.

#### 7. resourcedetectionprocessor

**v0.165.0 Changes:**
- The `system` and `docker` detectors set `os.type` to `zos` instead of `z_os` on z/OS ([#52804](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/52804))
- The `system` detector adds the `os.build_id` resource attribute and deprecates `os.build.id` ([#52806](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/52806))

**v0.166.0 Changes:**
- The `gcp` detector emits `gcp.cloud_run.job.task_index` as an integer instead of a string ([#52807](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/52807))

#### 8. routingconnector

**v0.165.0 Changes:**
- The `connector.routing.defaultErrorModeIgnore` feature gate was promoted to stable ([#52352](https://github.com/open-telemetry/opentelemetry-collector-contrib/pull/52352))
- **Impact**: The default `error_mode: ignore` is now permanent.
- **Action Required**: If you opted out via the feature gate, remove it and set `error_mode: propagate` explicitly in the connector config instead.

### sacloudexporter

No configuration change is needed. The changes that affect it are described above: the metric label name change, the retry delay of metrics, and the double batching warning.

### Components with No Breaking Changes

The following bundled components had no component-specific breaking changes in the v0.162.0–v0.166.0 range:
`otlpreceiver`, `otlpexporter`, `otlphttpexporter`, `debugexporter`, `batchprocessor`, `memorylimiterprocessor`, `awss3exporter`, `fileexporter`, `loadbalancingexporter`, `healthcheckextension`, `filestorage`, `attributesprocessor`, `groupbyattrsprocessor`, `resourceprocessor`, `deltatocumulativeprocessor`, `tailsamplingprocessor`, `hostmetricsreceiver`, `filelogreceiver`, `fluentforwardreceiver`, `syslogreceiver`, `journaldreceiver`, `windowseventlogreceiver`, `spanmetricsconnector`, `countconnector`.

> Note: `mackerelotlpexporter` is a third-party component (`github.com/mackerelio/opentelemetry-collector-mackerel`) with its own versioning. Its version is unchanged from v0.8.x (v0.15.1).

## Configuration Validation

After updating configurations, validate them using:

```bash
./sacloud-otel-collector validate --config your-config.yaml
```

The removed flat keepalive fields, the `elasticsearch` and `kafka` option removals and the OTTL changes are the ones most likely to make an existing `config.yaml` fail validation outright. The metric label name change and the `jaeger` receiver attribute change are the ones most likely to silently change behavior while still validating.

## References

### OpenTelemetry Collector Releases
- [v0.162.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.162.0)
- [v0.163.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.163.0)
- [v0.164.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.164.0)
- [v0.165.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.165.0)
- [v0.166.0](https://github.com/open-telemetry/opentelemetry-collector/releases/tag/v0.166.0)

### OpenTelemetry Collector Contrib Releases
- [v0.162.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.162.0)
- [v0.163.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.163.0)
- [v0.164.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.164.0)
- [v0.165.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.165.0)
- [v0.166.0](https://github.com/open-telemetry/opentelemetry-collector-contrib/releases/tag/v0.166.0)

## Support

If you encounter issues during the upgrade:

1. Check component-specific documentation at [OpenTelemetry Docs](https://opentelemetry.io/docs/)
2. Review the full changelog for each version
3. Report issues at the [sacloud-otel-collector repository](https://github.com/sacloud/sacloud-otel-collector/issues)
