//go:build windows

package e2e_test

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestWindowsEventLogToSakumock(t *testing.T) {
	collector, sakumock := findBinaries(t)

	dataPlaneAddr := freeLoopbackAddr(t)
	healthCheckAddr := freeLoopbackAddr(t)

	dumpDir := t.TempDir()
	startProcess(t, "sakumock", sakumock,
		"--enable-data-plane",
		"--data-plane-addr", dataPlaneAddr,
		"--data-plane-dump-dir", dumpDir,
	)
	waitListen(t, dataPlaneAddr, 30*time.Second)

	const eventSource = "SacloudOtelE2E"
	createEventSource(t, eventSource)

	cfg := fmt.Sprintf(`
receivers:
  windowseventlog:
    channel: Application
    start_at: end
exporters:
  sacloud:
    logs:
      endpoint: http://%s
      token: log-dummy
extensions:
  health_check:
    endpoint: %s
service:
  telemetry:
    metrics:
      level: none
  extensions: [health_check]
  pipelines:
    logs:
      receivers: [windowseventlog]
      exporters: [sacloud]
`, dataPlaneAddr, healthCheckAddr)
	cfgFile := writeConfigFile(t, cfg)

	collectorLog := startProcess(t, "collector", collector, "--config", cfgFile)
	waitListen(t, healthCheckAddr, 60*time.Second)

	marker := "sacloud-otel-collector e2e windowseventlog " + randomHex(t)
	writeEventLog(t, eventSource, marker)

	if !waitForDumpContaining(dumpDir, "otlp-logs-", marker, 60*time.Second) {
		t.Errorf("no otlp-logs-* dump containing %q in %s; collector log:\n%s",
			marker, dumpDir, readFile(collectorLog))
	}
}

func createEventSource(t *testing.T, source string) {
	t.Helper()
	// New-EventLog fails if the source already exists; ignore errors.
	_ = exec.Command("powershell", "-Command",
		fmt.Sprintf(`New-EventLog -LogName Application -Source %q -ErrorAction SilentlyContinue`, source),
	).Run()
}

func writeEventLog(t *testing.T, source, message string) {
	t.Helper()
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`Write-EventLog -LogName Application -Source %q -EventId 1000 -EntryType Information -Message %q`,
			source, message),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Write-EventLog failed: %v\n%s", err, out)
	}
}
