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
)

// newMetricsExporter creates a new metrics exporter using prometheusremotewriteexporter.
func newMetricsExporter(ctx context.Context, set exporter.Settings, cfg *Config) (exporter.Metrics, error) {
	factory := prometheusremotewriteexporter.NewFactory()
	defaultCfg := factory.CreateDefaultConfig()
	prwCfg, ok := defaultCfg.(*prometheusremotewriteexporter.Config)
	if !ok {
		return nil, fmt.Errorf("failed to cast to prometheusremotewriteexporter.Config")
	}

	// Configure endpoint
	prwCfg.ClientConfig.Endpoint = cfg.MetricsEndpointURL()

	// Configure timeout
	prwCfg.ClientConfig.Timeout = cfg.GetTimeout()

	// Configure authentication header
	prwCfg.ClientConfig.Headers.Set("Authorization", configopaque.String("Bearer "+string(cfg.Metrics.Token)))

	// Enable compression (snappy is required by Prometheus remote write protocol)
	prwCfg.ClientConfig.Compression = configcompression.TypeSnappy

	// Enable resource to telemetry conversion
	prwCfg.ResourceToTelemetrySettings.Enabled = true

	// Apply retry configuration
	prwCfg.BackOffConfig = cfg.GetRetryConfig()

	// Apply remote write queue configuration
	prwCfg.RemoteWriteQueue.Enabled = true
	prwCfg.RemoteWriteQueue.QueueSize = defaultRemoteWriteQueueSize
	prwCfg.RemoteWriteQueue.NumConsumers = defaultRemoteWriteNumConsumers

	// Apply remote write batch configuration
	prwCfg.MaxBatchSizeBytes = defaultRemoteWriteBatchSizeBytes

	// Create new settings with the correct component type
	prwSet := exporter.Settings{
		ID:                component.NewIDWithName(factory.Type(), set.ID.Name()),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	return factory.CreateMetrics(ctx, prwSet, prwCfg)
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
	otlpCfg.QueueConfig = configoptional.Some(defaultSendingQueueConfig())

	// Create new settings with the correct component type
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
	otlpCfg.QueueConfig = configoptional.Some(defaultSendingQueueConfig())

	// Create new settings with the correct component type
	otlpSet := exporter.Settings{
		ID:                component.NewIDWithName(factory.Type(), set.ID.Name()),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	return factory.CreateTraces(ctx, otlpSet, otlpCfg)
}
