package selfmetricsreceiver

import (
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "default config",
			cfg: Config{
				Endpoint:           defaultEndpoint,
				CollectionInterval: defaultCollectionInterval,
			},
			wantErr: false,
		},
		{
			name: "empty endpoint",
			cfg: Config{
				CollectionInterval: defaultCollectionInterval,
			},
			wantErr: true,
		},
		{
			name: "endpoint without port",
			cfg: Config{
				Endpoint:           "127.0.0.1",
				CollectionInterval: defaultCollectionInterval,
			},
			wantErr: true,
		},
		{
			name: "zero collection interval",
			cfg: Config{
				Endpoint: defaultEndpoint,
			},
			wantErr: true,
		},
		{
			name: "negative collection interval",
			cfg: Config{
				Endpoint:           defaultEndpoint,
				CollectionInterval: -time.Second,
			},
			wantErr: true,
		},
		{
			name: "custom endpoint and interval",
			cfg: Config{
				Endpoint:           "localhost:18888",
				CollectionInterval: 15 * time.Second,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
