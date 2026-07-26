package selfmetricsreceiver

import (
	"context"
	"testing"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	if factory.Type() != component.MustNewType(typeStr) {
		t.Errorf("factory.Type() = %v, want %v", factory.Type(), typeStr)
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	cfg, ok := createDefaultConfig().(*Config)
	if !ok {
		t.Fatal("createDefaultConfig() did not return *Config")
	}
	if cfg.Endpoint != defaultEndpoint {
		t.Errorf("Endpoint = %v, want %v", cfg.Endpoint, defaultEndpoint)
	}
	if cfg.CollectionInterval != defaultCollectionInterval {
		t.Errorf("CollectionInterval = %v, want %v", cfg.CollectionInterval, defaultCollectionInterval)
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config must be valid, got error: %v", err)
	}
}

func TestCreateMetricsReceiver(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	set := receivertest.NewNopSettings(factory.Type())
	recv, err := factory.CreateMetrics(context.Background(), set, cfg, consumertest.NewNop())
	if err != nil {
		t.Fatalf("CreateMetrics() error = %v", err)
	}
	if recv == nil {
		t.Fatal("CreateMetrics() returned nil receiver")
	}
}
