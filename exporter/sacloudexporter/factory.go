package sacloudexporter

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
)

const (
	// typeStr is the type identifier for the exporter.
	typeStr = "sacloud"
)

// NewFactory creates a new factory for the SAKURA Cloud exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		component.MustNewType(typeStr),
		createDefaultConfig,
		exporter.WithMetrics(createMetricsExporter, component.StabilityLevelDevelopment),
		exporter.WithLogs(createLogsExporter, component.StabilityLevelDevelopment),
		exporter.WithTraces(createTracesExporter, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createMetricsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Metrics, error) {
	oCfg := cfg.(*Config)
	if oCfg.Metrics.Endpoint == "" {
		return nil, fmt.Errorf("metrics.endpoint is not set")
	}
	return newMetricsExporter(ctx, set, oCfg)
}

func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	oCfg := cfg.(*Config)
	if oCfg.Logs.Endpoint == "" {
		return nil, fmt.Errorf("logs.endpoint is not set")
	}
	return newLogsExporter(ctx, set, oCfg)
}

func createTracesExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Traces, error) {
	oCfg := cfg.(*Config)
	if oCfg.Traces.Endpoint == "" {
		return nil, fmt.Errorf("traces.endpoint is not set")
	}
	return newTracesExporter(ctx, set, oCfg)
}
