package poller

import (
	"net/http"
	"strings"
	"time"

	"github.com/davidnewhall/unifi-poller/pkg/metrics"
	"github.com/davidnewhall/unifi-poller/pkg/promunifi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const oneDecimalPoint = 10

// RunPrometheus starts the web server and registers the collector.
func (u *UnifiPoller) RunPrometheus() error {
	u.Logf("Exporting Measurements at https://%s/metrics for Prometheus", u.Config.HTTPListen)
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(promunifi.NewUnifiCollector(promunifi.UnifiCollectorCnfg{
		Namespace:    strings.Replace(u.Config.Namespace, "-", "", -1),
		CollectFn:    u.ExportMetrics,
		LoggingFn:    u.LogExportReport,
		ReportErrors: true, // XXX: Does this need to be configurable?
	}))
	return http.ListenAndServe(u.Config.HTTPListen, nil)
}

// ExportMetrics updates the internal metrics provided via
// HTTP at /metrics for prometheus collection.
// This is run by Prometheus as CollectFn.
func (u *UnifiPoller) ExportMetrics() (*metrics.Metrics, error) {
	m, err := u.CollectMetrics()
	if err != nil {
		u.LogErrorf("collecting metrics: %v", err)
		u.Logf("Re-authenticating to UniFi Controller")
		if err := u.Unifi.Login(); err != nil {
			u.LogError(err, "re-authenticating")
			return nil, err
		}

		if m, err = u.CollectMetrics(); err != nil {
			u.LogErrorf("collecting metrics: %v", err)
			return nil, err
		}
	}

	u.AugmentMetrics(m)
	return m, nil
}

// LogExportReport is called after prometheus exports metrics.
// This is run by Prometheus as LoggingFn
func (u *UnifiPoller) LogExportReport(report *promunifi.Report) {
	m := report.Metrics
	u.Logf("UniFi Measurements Exported. Site: %d, Client: %d, "+
		"UAP: %d, USG/UDM: %d, USW: %d, Descs: %d, "+
		"Metrics: %d, Errs: %d, 0s: %d, Reqs/Total: %v / %v",
		len(m.Sites), len(m.Clients), len(m.UAPs), len(m.UDMs)+len(m.USGs), len(m.USWs),
		report.Descs, report.Total, report.Errors, report.Zeros,
		report.Fetch.Round(time.Millisecond/oneDecimalPoint),
		report.Elapsed.Round(time.Millisecond/oneDecimalPoint))
}
