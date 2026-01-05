package sacloudexporter

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// Default endpoint URL patterns for SAKURA Cloud Monitoring Suite
	metricsEndpointPattern = "https://%s.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write"
	logsEndpointPattern    = "https://%s.logs.monitoring.global.api.sacloud.jp"
	tracesEndpointPattern  = "https://%s.traces.monitoring.global.api.sacloud.jp"

	// Default queue settings for logs/traces (QueueBatchConfig)
	// These defaults ensure safe operation within SAKURA Cloud Monitoring Suite limits:
	// - 5 MB per request limit (using 4 MiB to be safe)
	// - 50 requests/second, 1000 lines/second rate limits
	defaultSendingQueueSize    = 10 * 1024 * 1024 // 10 MiB buffer
	defaultSendingNumConsumers = 2
	defaultBatchMaxSize        = 4 * 1024 * 1024  // 4 MiB per request (under 5 MB limit)
	defaultBatchFlushTimeout   = 10 * time.Second // batch logs for efficiency

	// Default queue settings for metrics (RemoteWriteQueue)
	defaultRemoteWriteQueueSize      = 10000
	defaultRemoteWriteNumConsumers   = 2
	defaultRemoteWriteBatchSizeBytes = 4 * 1024 * 1024 // 4 MiB per request (under 5 MB limit)

	// Default timeout for HTTP requests
	defaultTimeout = 30 * time.Second
)

// Config defines configuration for the SAKURA Cloud exporter.
type Config struct {
	// TimeoutConfig for timeout. Default is 30 seconds.
	exporterhelper.TimeoutConfig `mapstructure:",squash"`

	// BackOffConfig for retry on failure. Default is enabled with exponential backoff.
	configretry.BackOffConfig `mapstructure:"retry_on_failure"`

	// Metrics configuration for SAKURA Cloud Monitoring Suite metrics storage.
	Metrics MetricsEndpointConfig `mapstructure:"metrics"`

	// Logs configuration for SAKURA Cloud Monitoring Suite logs storage.
	Logs EndpointConfig `mapstructure:"logs"`

	// Traces configuration for SAKURA Cloud Monitoring Suite traces storage.
	Traces EndpointConfig `mapstructure:"traces"`
}

// MetricsEndpointConfig defines configuration for metrics signal.
type MetricsEndpointConfig struct {
	// Endpoint can be either:
	// - An endpoint identifier from SAKURA Cloud control panel (e.g., "123456789012")
	// - A full URL (e.g., "https://123456789012.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write")
	// If only an identifier is provided, it will be expanded to full URL automatically.
	Endpoint string `mapstructure:"endpoint"`

	// Token is the Bearer token for authentication.
	Token configopaque.String `mapstructure:"token"`
}

// EndpointConfig defines configuration for logs/traces signals.
type EndpointConfig struct {
	// Endpoint can be either:
	// - An endpoint identifier from SAKURA Cloud control panel (e.g., "123456789012")
	// - A full URL (e.g., "https://123456789012.logs.monitoring.global.api.sacloud.jp")
	// If only an identifier is provided, it will be expanded to full URL automatically.
	Endpoint string `mapstructure:"endpoint"`

	// Token is the Bearer token for authentication.
	Token configopaque.String `mapstructure:"token"`
}

// Validate checks if the configuration is valid.
func (cfg *Config) Validate() error {
	var errs []error

	// At least one endpoint must be configured
	if cfg.Metrics.Endpoint == "" && cfg.Logs.Endpoint == "" && cfg.Traces.Endpoint == "" {
		errs = append(errs, errors.New("at least one of metrics, logs, or traces endpoint must be configured"))
	}

	// Validate metrics config if endpoint is set
	if cfg.Metrics.Endpoint != "" && cfg.Metrics.Token == "" {
		errs = append(errs, errors.New("metrics.token is required when metrics.endpoint is set"))
	}

	// Validate logs config if endpoint is set
	if cfg.Logs.Endpoint != "" && cfg.Logs.Token == "" {
		errs = append(errs, errors.New("logs.token is required when logs.endpoint is set"))
	}

	// Validate traces config if endpoint is set
	if cfg.Traces.Endpoint != "" && cfg.Traces.Token == "" {
		errs = append(errs, errors.New("traces.token is required when traces.endpoint is set"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// isFullURL returns true if the endpoint is a full URL (starts with https://).
func isFullURL(endpoint string) bool {
	return strings.HasPrefix(endpoint, "https://")
}

// MetricsEndpointURL returns the full URL for metrics endpoint.
func (cfg *Config) MetricsEndpointURL() string {
	if cfg.Metrics.Endpoint == "" {
		return ""
	}
	if isFullURL(cfg.Metrics.Endpoint) {
		// Full URL provided, use as-is
		return cfg.Metrics.Endpoint
	}
	// ID provided, expand to full URL
	return fmt.Sprintf(metricsEndpointPattern, cfg.Metrics.Endpoint)
}

// LogsEndpointURL returns the full URL for logs endpoint.
func (cfg *Config) LogsEndpointURL() string {
	if cfg.Logs.Endpoint == "" {
		return ""
	}
	if isFullURL(cfg.Logs.Endpoint) {
		// Full URL provided, use as-is
		return cfg.Logs.Endpoint
	}
	// ID provided, expand to full URL
	return fmt.Sprintf(logsEndpointPattern, cfg.Logs.Endpoint)
}

// TracesEndpointURL returns the full URL for traces endpoint.
func (cfg *Config) TracesEndpointURL() string {
	if cfg.Traces.Endpoint == "" {
		return ""
	}
	if isFullURL(cfg.Traces.Endpoint) {
		// Full URL provided, use as-is
		return cfg.Traces.Endpoint
	}
	// ID provided, expand to full URL
	return fmt.Sprintf(tracesEndpointPattern, cfg.Traces.Endpoint)
}

// defaultSendingQueueConfig returns the default QueueBatchConfig for logs/traces.
// This configuration ensures safe operation within SAKURA Cloud Monitoring Suite limits.
func defaultSendingQueueConfig() exporterhelper.QueueBatchConfig {
	return exporterhelper.QueueBatchConfig{
		Sizer:        exporterhelper.RequestSizerTypeBytes,
		QueueSize:    defaultSendingQueueSize,
		NumConsumers: defaultSendingNumConsumers,
		Batch: configoptional.Some(exporterhelper.BatchConfig{
			FlushTimeout: defaultBatchFlushTimeout,
			Sizer:        exporterhelper.RequestSizerTypeBytes,
			MaxSize:      defaultBatchMaxSize,
		}),
	}
}

// GetTimeout returns the configured timeout or the default if not set.
func (cfg *Config) GetTimeout() time.Duration {
	if cfg.TimeoutConfig.Timeout == 0 {
		return defaultTimeout
	}
	return cfg.TimeoutConfig.Timeout
}

// GetRetryConfig returns the configured retry config or the default if not set.
func (cfg *Config) GetRetryConfig() configretry.BackOffConfig {
	if isZeroBackOffConfig(cfg.BackOffConfig) {
		return configretry.NewDefaultBackOffConfig()
	}
	return cfg.BackOffConfig
}

// isZeroBackOffConfig returns true if the BackOffConfig is not configured (zero value).
func isZeroBackOffConfig(cfg configretry.BackOffConfig) bool {
	return !cfg.Enabled && cfg.InitialInterval == 0 && cfg.MaxInterval == 0 && cfg.MaxElapsedTime == 0
}
