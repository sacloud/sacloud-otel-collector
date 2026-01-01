package sacloudexporter

import (
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/config/configopaque"
)

const (
	// Default endpoint URL patterns for SAKURA Cloud Monitoring Suite
	metricsEndpointPattern = "https://%s.metrics.monitoring.global.api.sacloud.jp/prometheus/api/v1/write"
	logsEndpointPattern    = "https://%s.logs.monitoring.global.api.sacloud.jp"
	tracesEndpointPattern  = "https://%s.traces.monitoring.global.api.sacloud.jp"
)

// Config defines configuration for the SAKURA Cloud exporter.
type Config struct {
	// Metrics configuration for SAKURA Cloud Monitoring Suite metrics storage.
	Metrics EndpointConfig `mapstructure:"metrics"`

	// Logs configuration for SAKURA Cloud Monitoring Suite logs storage.
	Logs EndpointConfig `mapstructure:"logs"`

	// Traces configuration for SAKURA Cloud Monitoring Suite traces storage.
	Traces EndpointConfig `mapstructure:"traces"`
}

// EndpointConfig defines configuration for each signal type.
type EndpointConfig struct {
	// Endpoint can be either:
	// - An endpoint identifier from SAKURA Cloud control panel (e.g., "123456789012")
	// - A full FQDN (e.g., "123456789012.logs.monitoring.global.api.sacloud.jp")
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

// isFQDN returns true if the endpoint looks like a FQDN (contains a dot).
func isFQDN(endpoint string) bool {
	return strings.Contains(endpoint, ".")
}

// MetricsEndpointURL returns the full URL for metrics endpoint.
func (cfg *Config) MetricsEndpointURL() string {
	if cfg.Metrics.Endpoint == "" {
		return ""
	}
	if isFQDN(cfg.Metrics.Endpoint) {
		// FQDN provided, add https:// prefix and metrics path
		return "https://" + cfg.Metrics.Endpoint + "/prometheus/api/v1/write"
	}
	// ID provided, expand to full URL
	return fmt.Sprintf(metricsEndpointPattern, cfg.Metrics.Endpoint)
}

// LogsEndpointURL returns the full URL for logs endpoint.
func (cfg *Config) LogsEndpointURL() string {
	if cfg.Logs.Endpoint == "" {
		return ""
	}
	if isFQDN(cfg.Logs.Endpoint) {
		// FQDN provided, add https:// prefix
		return "https://" + cfg.Logs.Endpoint
	}
	// ID provided, expand to full URL
	return fmt.Sprintf(logsEndpointPattern, cfg.Logs.Endpoint)
}

// TracesEndpointURL returns the full URL for traces endpoint.
func (cfg *Config) TracesEndpointURL() string {
	if cfg.Traces.Endpoint == "" {
		return ""
	}
	if isFQDN(cfg.Traces.Endpoint) {
		// FQDN provided, add https:// prefix
		return "https://" + cfg.Traces.Endpoint
	}
	// ID provided, expand to full URL
	return fmt.Sprintf(tracesEndpointPattern, cfg.Traces.Endpoint)
}
