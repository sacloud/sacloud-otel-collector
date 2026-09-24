package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocURL(t *testing.T) {
	tests := []struct {
		name string
		m    Module
		want string
	}{
		{
			name: "core",
			m:    Module{GoMod: "go.opentelemetry.io/collector/receiver/otlpreceiver v0.161.0"},
			want: "https://github.com/open-telemetry/opentelemetry-collector/tree/receiver/otlpreceiver/v0.161.0/receiver/otlpreceiver",
		},
		{
			name: "contrib nested",
			m:    Module{GoMod: "github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage v0.161.0"},
			want: "https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/extension/storage/filestorage/v0.161.0/extension/storage/filestorage",
		},
		{
			name: "local",
			m:    Module{GoMod: "github.com/sacloud/sacloud-otel-collector/receiver/selfmetricsreceiver v0.0.0", Path: "./receiver/selfmetricsreceiver"},
			want: "receiver/selfmetricsreceiver/README.md",
		},
		{
			name: "explicit doc",
			m:    Module{GoMod: "github.com/sacloud/sacloud-otel-collector/exporter/sacloudexporter v0.0.0", Path: "./exporter/sacloudexporter", Doc: "#sakuracloud-monitoring-suite"},
			want: "#sakuracloud-monitoring-suite",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.docURL()
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestDocURLUnknownHost(t *testing.T) {
	m := Module{GoMod: "example.com/foo/bar v1.0.0"}
	if _, err := m.docURL(); err == nil {
		t.Error("expected error for unknown host")
	}
}

func TestRenderTablesRequiresMetadata(t *testing.T) {
	cfg := &BuilderConfig{
		Receivers: []Module{{GoMod: "go.opentelemetry.io/collector/receiver/otlpreceiver v0.161.0"}},
	}
	if _, err := renderTables(cfg); err == nil {
		t.Error("expected error for missing type and description")
	}
}

const testConfig = `
exporters:
  - gomod: go.opentelemetry.io/collector/exporter/debugexporter v0.161.0
    type: debug
    description: Debug exporter
receivers:
  - gomod: go.opentelemetry.io/collector/receiver/otlpreceiver v0.161.0
    type: otlp
    description: OpenTelemetry Protocol receiver
providers:
  - gomod: go.opentelemetry.io/collector/confmap/provider/envprovider v1.67.0
`

const testReadme = `# Title

` + beginMarker + `
old content
` + endMarker + `

## Next
`

const wantReadme = `# Title

` + beginMarker + `

### Receivers

| Component | Description | Documentation |
|-----------|-------------|---------------|
| otlp | OpenTelemetry Protocol receiver | [Documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/receiver/otlpreceiver/v0.161.0/receiver/otlpreceiver) |

### Exporters

| Component | Description | Documentation |
|-----------|-------------|---------------|
| debug | Debug exporter | [Documentation](https://github.com/open-telemetry/opentelemetry-collector/tree/exporter/debugexporter/v0.161.0/exporter/debugexporter) |

` + endMarker + `

## Next
`

func TestRun(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "builder-config.yaml")
	readmePath := filepath.Join(dir, "README.md")
	if err := os.WriteFile(configPath, []byte(testConfig), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(readmePath, []byte(testReadme), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := run(configPath, readmePath, true); err == nil || !strings.Contains(err.Error(), "out of date") {
		t.Fatalf("check should fail on stale README: %v", err)
	}
	if err := run(configPath, readmePath, false); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != wantReadme {
		t.Errorf("unexpected README:\n%s", got)
	}
	if err := run(configPath, readmePath, true); err != nil {
		t.Errorf("check should pass after update: %v", err)
	}
}

func TestReplaceBetweenMarkersMissing(t *testing.T) {
	if _, err := replaceBetweenMarkers("no markers", "x"); err == nil {
		t.Error("expected error for missing markers")
	}
}
