package selfmetricsreceiver

import (
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	// defaultEndpoint is the default address of the collector's internal
	// telemetry Prometheus endpoint (service::telemetry::metrics).
	defaultEndpoint = "127.0.0.1:8888"

	// defaultCollectionInterval is the default scrape interval.
	defaultCollectionInterval = 60 * time.Second
)

// Config defines configuration for the selfmetrics receiver.
type Config struct {
	// Endpoint is the host:port of the collector's internal telemetry
	// Prometheus endpoint. Default is "127.0.0.1:8888".
	// Change this only if service::telemetry::metrics is configured to
	// expose the endpoint on a different address.
	Endpoint string `mapstructure:"endpoint"`

	// CollectionInterval is the interval at which internal metrics are
	// scraped. Default is 60 seconds.
	CollectionInterval time.Duration `mapstructure:"collection_interval"`
}

// Validate checks if the configuration is valid.
func (cfg *Config) Validate() error {
	var errs []error

	if cfg.Endpoint == "" {
		errs = append(errs, errors.New("endpoint must not be empty"))
	} else if _, _, err := net.SplitHostPort(cfg.Endpoint); err != nil {
		errs = append(errs, fmt.Errorf(`endpoint must be in "host:port" form: %w`, err))
	}

	if cfg.CollectionInterval <= 0 {
		errs = append(errs, errors.New("collection_interval must be positive"))
	}

	return errors.Join(errs...)
}
