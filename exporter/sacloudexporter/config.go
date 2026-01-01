package sacloudexporter

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configoptional"
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
	defaultRemoteWriteQueueSize    = 10000
	defaultRemoteWriteNumConsumers = 2
)

// Config defines configuration for the SAKURA Cloud exporter.
type Config struct {
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
	// - A full FQDN (e.g., "123456789012.metrics.monitoring.global.api.sacloud.jp")
	// If only an identifier is provided, it will be expanded to full URL automatically.
	Endpoint string `mapstructure:"endpoint"`

	// Token is the Bearer token for authentication.
	Token configopaque.String `mapstructure:"token"`

	// RemoteWriteQueue allows to configure the remote write queue.
	RemoteWriteQueue RemoteWriteQueue `mapstructure:"remote_write_queue"`
}

// RemoteWriteQueue allows to configure the remote write queue for metrics.
// This mirrors prometheusremotewriteexporter.RemoteWriteQueue.
type RemoteWriteQueue struct {
	// Enabled if false the queue is not enabled, the export requests
	// are executed synchronously.
	Enabled bool `mapstructure:"enabled"`

	// QueueSize is the maximum number of OTLP metric batches allowed
	// in the queue at a given time.
	QueueSize int `mapstructure:"queue_size"`

	// NumConsumers configures the number of workers used by
	// the collector to fan out remote write requests.
	NumConsumers int `mapstructure:"num_consumers"`
}

// EndpointConfig defines configuration for logs/traces signals.
type EndpointConfig struct {
	// Endpoint can be either:
	// - An endpoint identifier from SAKURA Cloud control panel (e.g., "123456789012")
	// - A full FQDN (e.g., "123456789012.logs.monitoring.global.api.sacloud.jp")
	// If only an identifier is provided, it will be expanded to full URL automatically.
	Endpoint string `mapstructure:"endpoint"`

	// Token is the Bearer token for authentication.
	Token configopaque.String `mapstructure:"token"`

	// SendingQueue defines configuration for queueing and batching.
	SendingQueue exporterhelper.QueueBatchConfig `mapstructure:"sending_queue"`
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

// DefaultSendingQueueConfig returns the default QueueBatchConfig for logs/traces.
// This configuration ensures safe operation within SAKURA Cloud Monitoring Suite limits.
func DefaultSendingQueueConfig() exporterhelper.QueueBatchConfig {
	return exporterhelper.QueueBatchConfig{
		Enabled:      true,
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

// DefaultRemoteWriteQueueConfig returns the default RemoteWriteQueue for metrics.
func DefaultRemoteWriteQueueConfig() RemoteWriteQueue {
	return RemoteWriteQueue{
		Enabled:      true,
		QueueSize:    defaultRemoteWriteQueueSize,
		NumConsumers: defaultRemoteWriteNumConsumers,
	}
}

// isZeroSendingQueue returns true if the SendingQueue is not configured (zero value).
func isZeroSendingQueue(q exporterhelper.QueueBatchConfig) bool {
	return !q.Enabled && q.QueueSize == 0 && q.NumConsumers == 0
}

// isZeroRemoteWriteQueue returns true if the RemoteWriteQueue is not configured (zero value).
func isZeroRemoteWriteQueue(q RemoteWriteQueue) bool {
	return !q.Enabled && q.QueueSize == 0 && q.NumConsumers == 0
}
