package sacloudexporter

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configcompression"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.uber.org/zap"
)

// newMetricsExporter creates a new metrics exporter using prometheusremotewriteexporter.
func newMetricsExporter(ctx context.Context, set exporter.Settings, cfg *Config) (exporter.Metrics, error) {
	factory := prometheusremotewriteexporter.NewFactory()
	prwCfg, err := metricsExporterConfig(factory, cfg)
	if err != nil {
		return nil, err
	}

	// The prometheusremotewrite exporter's queue has no storage extension
	// support, so sending_queue.storage cannot be applied to metrics.
	if cfg.SendingQueue.Storage != nil {
		if cfg.Logs.Endpoint == "" && cfg.Traces.Endpoint == "" {
			// Storage has no effect at all in a metrics-only setup;
			// this is most likely a misconfiguration.
			set.Logger.Warn("sending_queue.storage has no effect: it is not supported for metrics and no logs/traces endpoints are configured",
				zap.Stringer("storage", cfg.SendingQueue.Storage))
		} else {
			set.Logger.Info("sending_queue.storage applies to logs/traces only; metrics are not persisted",
				zap.Stringer("storage", cfg.SendingQueue.Storage))
		}
	}

	// Create new settings with the correct component type
	prwSet := exporter.Settings{
		ID:                component.NewIDWithName(factory.Type(), set.ID.Name()),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	return factory.CreateMetrics(ctx, prwSet, prwCfg)
}

// metricsExporterConfig builds the prometheusremotewriteexporter config.
func metricsExporterConfig(factory exporter.Factory, cfg *Config) (*prometheusremotewriteexporter.Config, error) {
	prwCfg, ok := factory.CreateDefaultConfig().(*prometheusremotewriteexporter.Config)
	if !ok {
		return nil, fmt.Errorf("failed to cast to prometheusremotewriteexporter.Config")
	}

	// HTTP client settings must go to the nested HTTP block. The exporter
	// only reads cfg.HTTP; the deprecated squashed ClientConfig is copied
	// into it only when the config is unmarshaled from YAML.
	prwCfg.HTTP.Endpoint = cfg.MetricsEndpointURL()
	prwCfg.HTTP.Timeout = cfg.GetTimeout()
	prwCfg.HTTP.Headers.Set("Authorization", configopaque.String("Bearer "+string(cfg.Metrics.Token)))
	// Snappy compression required by the Prometheus remote write protocol
	// is applied by the exporter itself, not by the HTTP client.

	// Disable the exporterhelper timeout (default 5s) that wraps the whole
	// push including retries, as the otlphttp exporter does. Each request is
	// bounded by HTTP.Timeout and retries by retry_on_failure.max_elapsed_time.
	prwCfg.TimeoutSettings.Timeout = 0

	// Convert all resource attributes to metric labels
	prwCfg.ResourceConstantLabels.Included = []string{"*"}

	// Apply retry configuration
	prwCfg.BackOffConfig = cfg.GetRetryConfig()

	// Apply remote write queue configuration
	prwCfg.RemoteWriteQueue.Enabled = true
	prwCfg.RemoteWriteQueue.QueueSize = defaultRemoteWriteQueueSize
	prwCfg.RemoteWriteQueue.NumConsumers = defaultRemoteWriteNumConsumers

	// Apply remote write batch configuration
	prwCfg.MaxBatchSizeBytes = defaultRemoteWriteBatchSizeBytes

	return prwCfg, nil
}

// newLogsExporter creates a new logs exporter using otlphttpexporter.
func newLogsExporter(ctx context.Context, set exporter.Settings, cfg *Config) (exporter.Logs, error) {
	factory := otlphttpexporter.NewFactory()
	defaultCfg := factory.CreateDefaultConfig()
	otlpCfg, ok := defaultCfg.(*otlphttpexporter.Config)
	if !ok {
		return nil, fmt.Errorf("failed to cast to otlphttpexporter.Config")
	}

	// Configure endpoint
	otlpCfg.ClientConfig.Endpoint = cfg.LogsEndpointURL()

	// Configure timeout
	otlpCfg.ClientConfig.Timeout = cfg.GetTimeout()

	// Configure authentication header
	otlpCfg.ClientConfig.Headers.Set("Authorization", configopaque.String("Bearer "+string(cfg.Logs.Token)))

	// Enable compression
	otlpCfg.ClientConfig.Compression = configcompression.TypeGzip

	// Apply retry configuration
	otlpCfg.RetryConfig = cfg.GetRetryConfig()

	// Apply sending queue configuration
	otlpCfg.QueueConfig = configoptional.Some(cfg.sendingQueueConfig())

	// Create new settings with the correct component type.
	// The rewritten ID is part of the persistent queue's storage key;
	// changing it would orphan data already queued in a storage extension.
	otlpSet := exporter.Settings{
		ID:                component.NewIDWithName(factory.Type(), set.ID.Name()),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	return factory.CreateLogs(ctx, otlpSet, otlpCfg)
}

// newTracesExporter creates a new traces exporter using otlphttpexporter.
func newTracesExporter(ctx context.Context, set exporter.Settings, cfg *Config) (exporter.Traces, error) {
	factory := otlphttpexporter.NewFactory()
	defaultCfg := factory.CreateDefaultConfig()
	otlpCfg, ok := defaultCfg.(*otlphttpexporter.Config)
	if !ok {
		return nil, fmt.Errorf("failed to cast to otlphttpexporter.Config")
	}

	// Configure endpoint
	otlpCfg.ClientConfig.Endpoint = cfg.TracesEndpointURL()

	// Configure timeout
	otlpCfg.ClientConfig.Timeout = cfg.GetTimeout()

	// Configure authentication header
	otlpCfg.ClientConfig.Headers.Set("Authorization", configopaque.String("Bearer "+string(cfg.Traces.Token)))

	// Enable compression
	otlpCfg.ClientConfig.Compression = configcompression.TypeGzip

	// Apply retry configuration
	otlpCfg.RetryConfig = cfg.GetRetryConfig()

	// Apply sending queue configuration
	otlpCfg.QueueConfig = configoptional.Some(cfg.sendingQueueConfig())

	// Create new settings with the correct component type.
	// The rewritten ID is part of the persistent queue's storage key;
	// changing it would orphan data already queued in a storage extension.
	otlpSet := exporter.Settings{
		ID:                component.NewIDWithName(factory.Type(), set.ID.Name()),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	return factory.CreateTraces(ctx, otlpSet, otlpCfg)
}
