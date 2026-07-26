package selfmetricsreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	// typeStr is the type identifier for the receiver.
	typeStr = "selfmetrics"
)

// NewFactory creates a new factory for the selfmetrics receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		component.MustNewType(typeStr),
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, component.StabilityLevelAlpha),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Endpoint:           defaultEndpoint,
		CollectionInterval: defaultCollectionInterval,
	}
}

func createMetricsReceiver(
	ctx context.Context,
	set receiver.Settings,
	cfg component.Config,
	next consumer.Metrics,
) (receiver.Metrics, error) {
	oCfg := cfg.(*Config)
	return newMetricsReceiver(ctx, set, oCfg, next)
}
