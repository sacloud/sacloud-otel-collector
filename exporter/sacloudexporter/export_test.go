package sacloudexporter

import (
	"context"
	"testing"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/exporter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewMetricsExporter_StorageLog(t *testing.T) {
	storageID := component.MustNewID("file_storage")

	tests := []struct {
		name      string
		cfg       Config
		wantLevel zapcore.Level
		wantLog   bool
	}{
		{
			name: "no storage logs nothing",
			cfg: Config{
				Metrics: MetricsEndpointConfig{Endpoint: "123456789012", Token: "token"},
			},
			wantLog: false,
		},
		{
			name: "storage with metrics only warns",
			cfg: Config{
				Metrics:      MetricsEndpointConfig{Endpoint: "123456789012", Token: "token"},
				SendingQueue: SendingQueueConfig{Storage: &storageID},
			},
			wantLog:   true,
			wantLevel: zapcore.WarnLevel,
		},
		{
			name: "storage with logs endpoint informs",
			cfg: Config{
				Metrics:      MetricsEndpointConfig{Endpoint: "123456789012", Token: "token"},
				Logs:         EndpointConfig{Endpoint: "123456789012", Token: "token"},
				SendingQueue: SendingQueueConfig{Storage: &storageID},
			},
			wantLog:   true,
			wantLevel: zapcore.InfoLevel,
		},
		{
			name: "storage with traces endpoint informs",
			cfg: Config{
				Metrics:      MetricsEndpointConfig{Endpoint: "123456789012", Token: "token"},
				Traces:       EndpointConfig{Endpoint: "123456789012", Token: "token"},
				SendingQueue: SendingQueueConfig{Storage: &storageID},
			},
			wantLog:   true,
			wantLevel: zapcore.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, observed := observer.New(zapcore.InfoLevel)
			set := exporter.Settings{
				ID:                component.NewID(component.MustNewType("sacloud")),
				TelemetrySettings: componenttest.NewNopTelemetrySettings(),
			}
			set.Logger = zap.New(core)

			if _, err := newMetricsExporter(context.Background(), set, &tt.cfg); err != nil {
				t.Fatalf("newMetricsExporter() error = %v", err)
			}

			logs := observed.FilterMessageSnippet("sending_queue.storage").All()
			if !tt.wantLog {
				if len(logs) != 0 {
					t.Fatalf("expected no storage log, got %v", logs)
				}
				return
			}
			if len(logs) != 1 {
				t.Fatalf("expected exactly one storage log, got %v", logs)
			}
			if logs[0].Level != tt.wantLevel {
				t.Errorf("storage log level = %v, want %v", logs[0].Level, tt.wantLevel)
			}
		})
	}
}
