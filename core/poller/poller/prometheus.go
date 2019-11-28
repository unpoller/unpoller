package poller

import (
	"fmt"
	"time"

	"github.com/davidnewhall/unifi-poller/metrics"
	"github.com/davidnewhall/unifi-poller/promunifi"
)

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
	idsMsg := ""
	if u.Config.CollectIDS {
		idsMsg = fmt.Sprintf(", IDS Events: %d, ", len(m.IDSList))
	}

	u.Logf("UniFi Measurements Exported. Sites: %d, Clients: %d, "+
		"Wireless APs: %d, Gateways: %d, Switches: %d%s, Descs: %d, "+
		"Metrics: %d, Errors: %d, Zeros: %d, Elapsed: %v",
		len(m.Sites), len(m.Clients), len(m.UAPs), len(m.UDMs)+len(m.USGs),
		len(m.USWs), idsMsg, report.Descs, report.Total, report.Errors,
		report.Zeros, report.Elapsed.Round(time.Millisecond))
}
