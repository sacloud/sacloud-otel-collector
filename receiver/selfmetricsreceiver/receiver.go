package selfmetricsreceiver

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

// scrapeJobName is used as the Prometheus job name, which becomes the
// service.name resource attribute of the collected metrics.
const scrapeJobName = "sacloud-otel-collector"

// newMetricsReceiver creates a new metrics receiver using prometheusreceiver
// configured to scrape the collector's own internal telemetry endpoint.
func newMetricsReceiver(ctx context.Context, set receiver.Settings, cfg *Config, next consumer.Metrics) (receiver.Metrics, error) {
	factory := prometheusreceiver.NewFactory()
	defaultCfg := factory.CreateDefaultConfig()
	promCfg, ok := defaultCfg.(*prometheusreceiver.Config)
	if !ok {
		return nil, fmt.Errorf("failed to cast to prometheusreceiver.Config")
	}

	// Build the scrape config via confmap so that prometheusreceiver's own
	// unmarshaling applies the Prometheus config defaults.
	conf := confmap.NewFromStringMap(map[string]any{
		"config": map[string]any{
			"scrape_configs": []any{
				map[string]any{
					"job_name":        scrapeJobName,
					"scrape_interval": cfg.CollectionInterval.String(),
					"static_configs": []any{
						map[string]any{
							"targets": []any{cfg.Endpoint},
						},
					},
				},
			},
		},
	})
	if err := conf.Unmarshal(promCfg); err != nil {
		return nil, fmt.Errorf("failed to build prometheus receiver config: %w", err)
	}

	// Create new settings with the correct component type
	promSet := receiver.Settings{
		ID:                component.NewIDWithName(factory.Type(), set.ID.Name()),
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	return factory.CreateMetrics(ctx, promSet, promCfg, next)
}
