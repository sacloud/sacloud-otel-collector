package sacloudexporter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
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

func TestNewMetricsExporter_NoDeprecatedResourceToTelemetry(t *testing.T) {
	core, observed := observer.New(zapcore.InfoLevel)
	set := exporter.Settings{
		ID:                component.NewID(component.MustNewType("sacloud")),
		TelemetrySettings: componenttest.NewNopTelemetrySettings(),
	}
	set.Logger = zap.New(core)

	cfg := &Config{
		Metrics: MetricsEndpointConfig{Endpoint: "123456789012", Token: "token"},
	}
	if _, err := newMetricsExporter(context.Background(), set, cfg); err != nil {
		t.Fatalf("newMetricsExporter() error = %v", err)
	}

	if logs := observed.FilterMessageSnippet("resource_to_telemetry_conversion").All(); len(logs) != 0 {
		t.Fatalf("expected no resource_to_telemetry_conversion deprecation log, got %v", logs)
	}
}

func TestMetricsExporterConfig(t *testing.T) {
	cfg := &Config{
		Metrics: MetricsEndpointConfig{Endpoint: "123456789012", Token: "token"},
	}
	prwCfg, err := metricsExporterConfig(prometheusremotewriteexporter.NewFactory(), cfg)
	if err != nil {
		t.Fatalf("metricsExporterConfig() error = %v", err)
	}

	// The exporter only reads the nested HTTP block.
	wantEndpoint := "https://123456789012.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write"
	if prwCfg.HTTP.Endpoint != wantEndpoint {
		t.Errorf("HTTP.Endpoint = %q, want %q", prwCfg.HTTP.Endpoint, wantEndpoint)
	}
	if prwCfg.HTTP.Timeout != defaultTimeout {
		t.Errorf("HTTP.Timeout = %v, want %v", prwCfg.HTTP.Timeout, defaultTimeout)
	}
	if v, ok := prwCfg.HTTP.Headers.Get("Authorization"); !ok || v != "Bearer token" {
		t.Errorf("HTTP.Headers[Authorization] = %q, want %q", v, "Bearer token")
	}
	if prwCfg.TimeoutSettings.Timeout != 0 {
		t.Errorf("TimeoutSettings.Timeout = %v, want 0", prwCfg.TimeoutSettings.Timeout)
	}
}

func TestMetricsExporterConfig_CustomTimeout(t *testing.T) {
	cfg := &Config{
		TimeoutConfig: exporterhelper.TimeoutConfig{Timeout: 2 * time.Second},
		Metrics:       MetricsEndpointConfig{Endpoint: "123456789012", Token: "token"},
	}
	prwCfg, err := metricsExporterConfig(prometheusremotewriteexporter.NewFactory(), cfg)
	if err != nil {
		t.Fatalf("metricsExporterConfig() error = %v", err)
	}

	if prwCfg.HTTP.Timeout != 2*time.Second {
		t.Errorf("HTTP.Timeout = %v, want %v", prwCfg.HTTP.Timeout, 2*time.Second)
	}
	if prwCfg.TimeoutSettings.Timeout != 0 {
		t.Errorf("TimeoutSettings.Timeout = %v, want 0", prwCfg.TimeoutSettings.Timeout)
	}
}

func TestMetricsExporter_SendsToEndpoint(t *testing.T) {
	gotAuth := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case gotAuth <- r.Header.Get("Authorization"):
		default:
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	set := exporter.Settings{
		ID:                component.NewID(component.MustNewType("sacloud")),
		TelemetrySettings: componenttest.NewNopTelemetrySettings(),
	}
	cfg := &Config{
		Metrics: MetricsEndpointConfig{Endpoint: srv.URL + "/prometheus/api/v1/write", Token: "token"},
	}
	exp, err := newMetricsExporter(context.Background(), set, cfg)
	if err != nil {
		t.Fatalf("newMetricsExporter() error = %v", err)
	}
	if err := exp.Start(context.Background(), componenttest.NewNopHost()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	defer func() {
		if err := exp.Shutdown(context.Background()); err != nil {
			t.Errorf("Shutdown() error = %v", err)
		}
	}()

	md := pmetric.NewMetrics()
	m := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
	m.SetName("test_gauge")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetDoubleValue(1)
	dp.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
	if err := exp.ConsumeMetrics(context.Background(), md); err != nil {
		t.Fatalf("ConsumeMetrics() error = %v", err)
	}

	select {
	case auth := <-gotAuth:
		if auth != "Bearer token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer token")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("no remote write request received by the endpoint")
	}
}
