package sacloudexporter

import (
	"testing"
	"time"

	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "all endpoints empty",
			cfg:     Config{},
			wantErr: true,
		},
		{
			name: "metrics endpoint without token",
			cfg: Config{
				Metrics: MetricsEndpointConfig{
					Endpoint: "123456789012",
				},
			},
			wantErr: true,
		},
		{
			name: "logs endpoint without token",
			cfg: Config{
				Logs: EndpointConfig{
					Endpoint: "123456789012",
				},
			},
			wantErr: true,
		},
		{
			name: "traces endpoint without token",
			cfg: Config{
				Traces: EndpointConfig{
					Endpoint: "123456789012",
				},
			},
			wantErr: true,
		},
		{
			name: "valid metrics config",
			cfg: Config{
				Metrics: MetricsEndpointConfig{
					Endpoint: "123456789012",
					Token:    "test-token",
				},
			},
			wantErr: false,
		},
		{
			name: "valid logs config",
			cfg: Config{
				Logs: EndpointConfig{
					Endpoint: "123456789012",
					Token:    "test-token",
				},
			},
			wantErr: false,
		},
		{
			name: "valid traces config",
			cfg: Config{
				Traces: EndpointConfig{
					Endpoint: "123456789012",
					Token:    "test-token",
				},
			},
			wantErr: false,
		},
		{
			name: "valid all signals config",
			cfg: Config{
				Metrics: MetricsEndpointConfig{
					Endpoint: "123456789012",
					Token:    "met-token",
				},
				Logs: EndpointConfig{
					Endpoint: "123456789012",
					Token:    "log-token",
				},
				Traces: EndpointConfig{
					Endpoint: "123456789012",
					Token:    "trc-token",
				},
			},
			wantErr: false,
		},
		{
			name: "valid full URL endpoint",
			cfg: Config{
				Metrics: MetricsEndpointConfig{
					Endpoint: "https://custom.endpoint.example.com",
					Token:    "test-token",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid endpoint - http URL",
			cfg: Config{
				Metrics: MetricsEndpointConfig{
					Endpoint: "http://example.com",
					Token:    "test-token",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid endpoint - bare hostname",
			cfg: Config{
				Logs: EndpointConfig{
					Endpoint: "example.com",
					Token:    "test-token",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid endpoint - alphanumeric mix",
			cfg: Config{
				Traces: EndpointConfig{
					Endpoint: "abc123",
					Token:    "test-token",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		wantErr  bool
	}{
		{"", false},
		{"123456789012", false},
		{"0", false},
		{"https://example.com", false},
		{"https://123456789012.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write", false},
		{"http://example.com", true},
		{"example.com", true},
		{"abc123", true},
		{"abc", true},
		{"123abc", true},
		{"12345 6789", true},
	}

	for _, tt := range tests {
		t.Run(tt.endpoint, func(t *testing.T) {
			err := validateEndpoint(tt.endpoint, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEndpoint(%q) error = %v, wantErr %v", tt.endpoint, err, tt.wantErr)
			}
		})
	}
}

func TestIsFullURL(t *testing.T) {
	tests := []struct {
		endpoint string
		want     bool
	}{
		{"https://example.com", true},
		{"https://123456789012.logs.monitoring.global.api.sacloud.jp", true},
		{"123456789012", false},
		{"example.com", false},
		{"http://example.com", false}, // only https is supported
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.endpoint, func(t *testing.T) {
			if got := isFullURL(tt.endpoint); got != tt.want {
				t.Errorf("isFullURL(%q) = %v, want %v", tt.endpoint, got, tt.want)
			}
		})
	}
}

func TestMetricsEndpointURL(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{
			name:     "empty endpoint",
			endpoint: "",
			want:     "",
		},
		{
			name:     "ID only",
			endpoint: "123456789012",
			want:     "https://123456789012.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write",
		},
		{
			name:     "full URL",
			endpoint: "https://123456789012.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write",
			want:     "https://123456789012.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Metrics: MetricsEndpointConfig{
					Endpoint: tt.endpoint,
				},
			}
			if got := cfg.MetricsEndpointURL(); got != tt.want {
				t.Errorf("MetricsEndpointURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogsEndpointURL(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{
			name:     "empty endpoint",
			endpoint: "",
			want:     "",
		},
		{
			name:     "ID only",
			endpoint: "123456789012",
			want:     "https://123456789012.logs.monitoring.global.api.sacloud.jp",
		},
		{
			name:     "full URL",
			endpoint: "https://123456789012.logs.monitoring.global.api.sacloud.jp",
			want:     "https://123456789012.logs.monitoring.global.api.sacloud.jp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Logs: EndpointConfig{
					Endpoint: tt.endpoint,
				},
			}
			if got := cfg.LogsEndpointURL(); got != tt.want {
				t.Errorf("LogsEndpointURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTracesEndpointURL(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		want     string
	}{
		{
			name:     "empty endpoint",
			endpoint: "",
			want:     "",
		},
		{
			name:     "ID only",
			endpoint: "123456789012",
			want:     "https://123456789012.traces.monitoring.global.api.sacloud.jp",
		},
		{
			name:     "full URL",
			endpoint: "https://123456789012.traces.monitoring.global.api.sacloud.jp",
			want:     "https://123456789012.traces.monitoring.global.api.sacloud.jp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Traces: EndpointConfig{
					Endpoint: tt.endpoint,
				},
			}
			if got := cfg.TracesEndpointURL(); got != tt.want {
				t.Errorf("TracesEndpointURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		want    time.Duration
	}{
		{
			name:    "zero value returns default",
			timeout: 0,
			want:    30 * time.Second,
		},
		{
			name:    "custom timeout",
			timeout: 60 * time.Second,
			want:    60 * time.Second,
		},
		{
			name:    "short timeout",
			timeout: 5 * time.Second,
			want:    5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				TimeoutConfig: exporterhelper.TimeoutConfig{Timeout: tt.timeout},
			}
			if got := cfg.GetTimeout(); got != tt.want {
				t.Errorf("Config.GetTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsZeroBackOffConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  configretry.BackOffConfig
		want bool
	}{
		{
			name: "zero value",
			cfg:  configretry.BackOffConfig{},
			want: true,
		},
		{
			name: "enabled only",
			cfg: configretry.BackOffConfig{
				Enabled: true,
			},
			want: false,
		},
		{
			name: "initial_interval only",
			cfg: configretry.BackOffConfig{
				InitialInterval: 5 * time.Second,
			},
			want: false,
		},
		{
			name: "max_interval only",
			cfg: configretry.BackOffConfig{
				MaxInterval: 30 * time.Second,
			},
			want: false,
		},
		{
			name: "max_elapsed_time only",
			cfg: configretry.BackOffConfig{
				MaxElapsedTime: 5 * time.Minute,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isZeroBackOffConfig(tt.cfg); got != tt.want {
				t.Errorf("isZeroBackOffConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetRetryConfig(t *testing.T) {
	defaultCfg := configretry.NewDefaultBackOffConfig()

	tests := []struct {
		name       string
		cfg        configretry.BackOffConfig
		wantEnable bool
	}{
		{
			name:       "zero value returns default (enabled)",
			cfg:        configretry.BackOffConfig{},
			wantEnable: true,
		},
		{
			name: "custom config with enabled=false",
			cfg: configretry.BackOffConfig{
				Enabled:         false,
				InitialInterval: 10 * time.Second, // non-zero to avoid being treated as zero config
			},
			wantEnable: false,
		},
		{
			name: "custom config with enabled=true",
			cfg: configretry.BackOffConfig{
				Enabled:         true,
				InitialInterval: 10 * time.Second,
			},
			wantEnable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				BackOffConfig: tt.cfg,
			}
			got := cfg.GetRetryConfig()
			if got.Enabled != tt.wantEnable {
				t.Errorf("Config.GetRetryConfig().Enabled = %v, want %v", got.Enabled, tt.wantEnable)
			}
			// For zero config, verify it returns the actual default values
			if isZeroBackOffConfig(tt.cfg) {
				if got.InitialInterval != defaultCfg.InitialInterval {
					t.Errorf("Config.GetRetryConfig().InitialInterval = %v, want %v", got.InitialInterval, defaultCfg.InitialInterval)
				}
				if got.MaxInterval != defaultCfg.MaxInterval {
					t.Errorf("Config.GetRetryConfig().MaxInterval = %v, want %v", got.MaxInterval, defaultCfg.MaxInterval)
				}
				if got.MaxElapsedTime != defaultCfg.MaxElapsedTime {
					t.Errorf("Config.GetRetryConfig().MaxElapsedTime = %v, want %v", got.MaxElapsedTime, defaultCfg.MaxElapsedTime)
				}
			}
		})
	}
}
