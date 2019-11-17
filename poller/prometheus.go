package poller

import (
	"fmt"
	"time"

	"github.com/davidnewhall/unifi-poller/metrics"
)

// ExportMetrics updates the internal metrics provided via
// HTTP at /metrics for prometheus collection. This is run by Prometheus.
func (u *UnifiPoller) ExportMetrics() *metrics.Metrics {
	if u.Config.ReAuth {
		u.LogDebugf("Re-authenticating to UniFi Controller")
		// Some users need to re-auth every interval because the cookie times out.
		if err := u.Unifi.Login(); err != nil {
			u.LogError(err, "re-authenticating")
			return nil
		}
	}
	u.LastCheck = time.Now()
	m, err := u.CollectMetrics()
	if err != nil {
		u.LogErrorf("collecting metrics: %v", err)
		return nil
	}
	u.AugmentMetrics(m)

	idsMsg := ""
	if u.Config.CollectIDS {
		idsMsg = fmt.Sprintf(", IDS Events: %d, ", len(m.IDSList))
	}
	u.Logf("UniFi Measurements Exported. Sites: %d, Clients: %d, "+
		"Wireless APs: %d, Gateways: %d, Switches: %d%s",
		len(m.Sites), len(m.Clients), len(m.UAPs),
		len(m.UDMs)+len(m.USGs), len(m.USWs), idsMsg)

	return m
}
