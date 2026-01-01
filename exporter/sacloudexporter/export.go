package sacloudexporter

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configopaque"
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

	// Configure authentication header
	if prwCfg.ClientConfig.Headers == nil {
		prwCfg.ClientConfig.Headers = make(map[string]configopaque.String)
	}
	prwCfg.ClientConfig.Headers["Authorization"] = configopaque.String("Bearer " + string(cfg.Metrics.Token))

	// Enable resource to telemetry conversion
	prwCfg.ResourceToTelemetrySettings.Enabled = true

	// Apply remote write queue configuration (use defaults if not explicitly configured)
	rwq := cfg.Metrics.RemoteWriteQueue
	if isZeroRemoteWriteQueue(rwq) {
		rwq = DefaultRemoteWriteQueueConfig()
	}
	prwCfg.RemoteWriteQueue.Enabled = rwq.Enabled
	if rwq.QueueSize > 0 {
		prwCfg.RemoteWriteQueue.QueueSize = rwq.QueueSize
	}
	if rwq.NumConsumers > 0 {
		prwCfg.RemoteWriteQueue.NumConsumers = rwq.NumConsumers
	}

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

	// Configure authentication header
	if otlpCfg.ClientConfig.Headers == nil {
		otlpCfg.ClientConfig.Headers = make(map[string]configopaque.String)
	}
	otlpCfg.ClientConfig.Headers["Authorization"] = configopaque.String("Bearer " + string(cfg.Logs.Token))

	// Apply sending queue configuration (use defaults if not explicitly configured)
	if isZeroSendingQueue(cfg.Logs.SendingQueue) {
		otlpCfg.QueueConfig = DefaultSendingQueueConfig()
	} else {
		otlpCfg.QueueConfig = cfg.Logs.SendingQueue
	}

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

	// Configure authentication header
	if otlpCfg.ClientConfig.Headers == nil {
		otlpCfg.ClientConfig.Headers = make(map[string]configopaque.String)
	}
	otlpCfg.ClientConfig.Headers["Authorization"] = configopaque.String("Bearer " + string(cfg.Traces.Token))

	// Apply sending queue configuration (use defaults if not explicitly configured)
	if isZeroSendingQueue(cfg.Traces.SendingQueue) {
		otlpCfg.QueueConfig = DefaultSendingQueueConfig()
	} else {
		otlpCfg.QueueConfig = cfg.Traces.SendingQueue
	}

	// Create new settings with the correct component type
	otlpSet := exporter.Settings{
		ID:                component.NewIDWithName(factory.Type(), set.ID.Name()),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	return factory.CreateTraces(ctx, otlpSet, otlpCfg)
}
