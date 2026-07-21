// Package e2e runs the built sacloud-otel-collector binary against the
// sakumock monitoring-suite data plane (https://github.com/sacloud/sakumock):
// a filelog receiver tails a file, the sacloud exporter forwards its lines as
// OTLP/HTTP logs to the mock, and the test asserts the mock's JSON dump
// contains what was written.
//
// It skips unless ../sacloud-otel-collector (built by `make`) exists and
// sakumock is on PATH (CI downloads the release binary; locally either grab
// one from https://github.com/sacloud/sakumock/releases or run
// `go install github.com/sacloud/sakumock/cmd/sakumock@latest`).
package e2e_test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestFilelogToSakumock(t *testing.T) {
	collector, sakumock := findBinaries(t)

	dataPlaneAddr := freeLoopbackAddr(t)
	healthCheckAddr := freeLoopbackAddr(t)

	dumpDir := t.TempDir()
	startProcess(t, "sakumock", sakumock, "monitoringsuite",
		"--enable-data-plane",
		"--data-plane-addr", dataPlaneAddr,
		"--data-plane-dump-dir", dumpDir,
	)
	waitListen(t, dataPlaneAddr, 30*time.Second)

	logFile := filepath.Join(t.TempDir(), "app.log")
	if err := os.WriteFile(logFile, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := fmt.Sprintf(`
receivers:
  filelog:
    include: [%q]
    start_at: beginning
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
      receivers: [filelog]
      exporters: [sacloud]
`, filepath.ToSlash(logFile), dataPlaneAddr, healthCheckAddr)
	cfgFile := writeConfigFile(t, cfg)

	collectorLog := startProcess(t, "collector", collector, "--config", cfgFile)
	waitListen(t, healthCheckAddr, 60*time.Second)

	marker := "sacloud-otel-collector e2e marker " + randomHex(t)
	appendLine(t, logFile, marker)

	if !waitForDumpContaining(dumpDir, "otlp-logs-", marker, 60*time.Second) {
		t.Errorf("no otlp-logs-* dump containing %q in %s; collector log:\n%s",
			marker, dumpDir, readFile(collectorLog))
	}
}

func findBinaries(t *testing.T) (collector, sakumock string) {
	t.Helper()
	collector, err := filepath.Abs(filepath.Join("..", collectorBinName()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(collector); err != nil {
		t.Skipf("%s not found (build it with `make`); skipping e2e", collector)
	}
	sakumock, err = exec.LookPath("sakumock")
	if err != nil {
		t.Skip("sakumock not found in PATH; skipping e2e")
	}
	return collector, sakumock
}

func writeConfigFile(t *testing.T, cfg string) string {
	t.Helper()
	cfgFile := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	return cfgFile
}

func collectorBinName() string {
	if runtime.GOOS == "windows" {
		return "sacloud-otel-collector.exe"
	}
	return "sacloud-otel-collector"
}

// startProcess runs a binary as a subprocess killed on test cleanup and
// returns the path of the file capturing its combined output.
func startProcess(t *testing.T, name, bin string, args ...string) string {
	t.Helper()
	logf, err := os.CreateTemp(t.TempDir(), name+"-*.log")
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, args...)
	cmd.Stdout, cmd.Stderr = logf, logf
	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", name, err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		_ = logf.Close()
	})
	return logf.Name()
}

func appendLine(t *testing.T, path, line string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		t.Fatal(err)
	}
}

func randomHex(t *testing.T) string {
	t.Helper()
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(b)
}

func readFile(path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}

// freeLoopbackAddr reserves a loopback port by binding and immediately
// closing a listener.
func freeLoopbackAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("allocate free port: %v", err)
	}
	addr := l.Addr().String()
	_ = l.Close()
	return addr
}

func waitListen(t *testing.T, addr string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if c, err := net.DialTimeout("tcp", addr, time.Second); err == nil {
			_ = c.Close()
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s to listen", addr)
}

// waitForDumpContaining polls dir until a file with the given name prefix
// contains the marker string.
func waitForDumpContaining(dir, prefix, marker string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), prefix) &&
				strings.Contains(readFile(filepath.Join(dir, e.Name())), marker) {
				return true
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
	return false
}
