package selfmetricsreceiver

import (
	"context"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestBuildPromConfig(t *testing.T) {
	cfg := &Config{
		Endpoint:           "127.0.0.1:8888",
		CollectionInterval: 15 * time.Second,
	}

	promCfg, err := buildPromConfig(prometheusreceiver.NewFactory(), cfg)
	if err != nil {
		t.Fatalf("buildPromConfig() error = %v", err)
	}

	scrapeConfigs := promCfg.PrometheusConfig.ScrapeConfigs
	if len(scrapeConfigs) != 1 {
		t.Fatalf("len(ScrapeConfigs) = %d, want 1", len(scrapeConfigs))
	}
	sc := scrapeConfigs[0]

	if sc.JobName != scrapeJobName {
		t.Errorf("JobName = %v, want %v", sc.JobName, scrapeJobName)
	}
	if got := time.Duration(sc.ScrapeInterval); got != cfg.CollectionInterval {
		t.Errorf("ScrapeInterval = %v, want %v", got, cfg.CollectionInterval)
	}
	if len(sc.ServiceDiscoveryConfigs) != 1 {
		t.Fatalf("len(ServiceDiscoveryConfigs) = %d, want 1", len(sc.ServiceDiscoveryConfigs))
	}
}

func TestFilteringConsumerDropsScrapeMetrics(t *testing.T) {
	md := pmetric.NewMetrics()
	sm := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty()
	for _, name := range []string{
		"up",
		"scrape_duration_seconds",
		"scrape_samples_scraped",
		"scrape_samples_post_metric_relabeling",
		"scrape_series_added",
		"scrape_body_size_bytes",
		"otelcol_process_uptime",
		"otelcol_exporter_sent_metric_points",
		"uptime",
		"my_scrape_metric",
	} {
		m := sm.Metrics().AppendEmpty()
		m.SetName(name)
		m.SetEmptyGauge().DataPoints().AppendEmpty().SetDoubleValue(1)
	}

	sink := &consumertest.MetricsSink{}
	fc := &filteringConsumer{next: sink}
	if err := fc.ConsumeMetrics(context.Background(), md); err != nil {
		t.Fatalf("ConsumeMetrics() error = %v", err)
	}

	got := map[string]bool{}
	for _, out := range sink.AllMetrics() {
		rms := out.ResourceMetrics()
		for i := 0; i < rms.Len(); i++ {
			sms := rms.At(i).ScopeMetrics()
			for j := 0; j < sms.Len(); j++ {
				ms := sms.At(j).Metrics()
				for k := 0; k < ms.Len(); k++ {
					got[ms.At(k).Name()] = true
				}
			}
		}
	}

	dropped := []string{
		"up",
		"scrape_duration_seconds",
		"scrape_samples_scraped",
		"scrape_samples_post_metric_relabeling",
		"scrape_series_added",
		"scrape_body_size_bytes",
	}
	for _, name := range dropped {
		if got[name] {
			t.Errorf("metric %q must be dropped", name)
		}
	}

	kept := []string{
		"otelcol_process_uptime",
		"otelcol_exporter_sent_metric_points",
		"uptime",
		"my_scrape_metric",
	}
	for _, name := range kept {
		if !got[name] {
			t.Errorf("metric %q must be kept", name)
		}
	}
}

func TestFilteringConsumerRemovesEmptyResources(t *testing.T) {
	md := pmetric.NewMetrics()
	sm := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty()
	m := sm.Metrics().AppendEmpty()
	m.SetName("up")
	m.SetEmptyGauge().DataPoints().AppendEmpty().SetDoubleValue(1)

	sink := &consumertest.MetricsSink{}
	fc := &filteringConsumer{next: sink}
	if err := fc.ConsumeMetrics(context.Background(), md); err != nil {
		t.Fatalf("ConsumeMetrics() error = %v", err)
	}

	for _, out := range sink.AllMetrics() {
		if out.ResourceMetrics().Len() != 0 {
			t.Errorf("ResourceMetrics().Len() = %d, want 0", out.ResourceMetrics().Len())
		}
	}
}
