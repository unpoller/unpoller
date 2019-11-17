package poller

import (
	"fmt"
	"time"

	"github.com/davidnewhall/unifi-poller/metrics"
)

// ExportMetrics updates the internal metrics provided via
// HTTP at /metrics for prometheus collection. This is run by Prometheus CollectFn.
func (u *UnifiPoller) ExportMetrics() (*metrics.Metrics, error) {
	if u.Config.ReAuth {
		u.LogDebugf("Re-authenticating to UniFi Controller")
		// Some users need to re-auth every interval because the cookie times out.
		if err := u.Unifi.Login(); err != nil {
			u.LogError(err, "re-authenticating")
			return nil, err
		}
	}
	u.LastCheck = time.Now()
	m, err := u.CollectMetrics()
	if err != nil {
		u.LogErrorf("collecting metrics: %v", err)
		return nil, err
	}
	u.AugmentMetrics(m)
	return m, nil
}

// LogExportReport is called after prometheus exports metrics. This is run by Prometheus as LoggingFn
func (u *UnifiPoller) LogExportReport(m *metrics.Metrics, count int64) {
	idsMsg := ""
	if u.Config.CollectIDS {
		idsMsg = fmt.Sprintf(", IDS Events: %d, ", len(m.IDSList))
	}
	u.Logf("UniFi Measurements Exported. Sites: %d, Clients: %d, "+
		"Wireless APs: %d, Gateways: %d, Switches: %d%s, Metrics: %d",
		len(m.Sites), len(m.Clients), len(m.UAPs),
		len(m.UDMs)+len(m.USGs), len(m.USWs), idsMsg, count)
}
