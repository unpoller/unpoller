package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// reportSite generates Unifi Sites' datapoints for Datadog.
// These points can be passed directly to Datadog.
func (u *DatadogUnifi) reportSite(r report, s *unifi.Site) {
	metricName := metricNamespace("subsystems")

	for _, h := range s.Health {
		tags := []string{
			tag("name", s.Name),
			tag("site_name", s.SiteName),
			tag("source", s.SourceName),
			tag("desc", s.Desc),
			tag("status", h.Status),
			tag("subsystem", h.Subsystem),
			tag("wan_ip", h.WanIP),
			tag("gw_name", h.GwName),
			tag("lan_ip", h.LanIP),
		}

		data := map[string]float64{
			"num_user":                 h.NumUser.Val,
			"num_guest":                h.NumGuest.Val,
			"num_iot":                  h.NumIot.Val,
			"tx_bytes_r":               h.TxBytesR.Val,
			"rx_bytes_r":               h.RxBytesR.Val,
			"num_ap":                   h.NumAp.Val,
			"num_adopted":              h.NumAdopted.Val,
			"num_disabled":             h.NumDisabled.Val,
			"num_disconnected":         h.NumDisconnected.Val,
			"num_pending":              h.NumPending.Val,
			"num_gw":                   h.NumGw.Val,
			"num_sta":                  h.NumSta.Val,
			"gw_cpu":                   h.GwSystemStats.CPU.Val,
			"gw_mem":                   h.GwSystemStats.Mem.Val,
			"gw_uptime":                h.GwSystemStats.Uptime.Val,
			"latency":                  h.Latency.Val,
			"uptime":                   h.Uptime.Val,
			"drops":                    h.Drops.Val,
			"xput_up":                  h.XputUp.Val,
			"xput_down":                h.XputDown.Val,
			"speedtest_ping":           h.SpeedtestPing.Val,
			"speedtest_lastrun":        h.SpeedtestLastrun.Val,
			"num_sw":                   h.NumSw.Val,
			"remote_user_num_active":   h.RemoteUserNumActive.Val,
			"remote_user_num_inactive": h.RemoteUserNumInactive.Val,
			"remote_user_rx_bytes":     h.RemoteUserRxBytes.Val,
			"remote_user_tx_bytes":     h.RemoteUserTxBytes.Val,
			"remote_user_rx_packets":   h.RemoteUserRxPackets.Val,
			"remote_user_tx_packets":   h.RemoteUserTxPackets.Val,
			"num_new_alarms":           s.NumNewAlarms.Val,
		}

		for name, value := range data {
			_ = r.reportGauge(metricName(name), value, tags)
		}
	}
}

func (u *DatadogUnifi) reportSiteDPI(r report, s *unifi.DPITable) {
	for _, dpi := range s.ByApp {
		metricName := metricNamespace("sitedpi")

		tags := []string{
			tag("category", unifi.DPICats.Get(dpi.Cat.Int())),
			tag("application", unifi.DPIApps.GetApp(dpi.Cat.Int(), dpi.App.Int())),
			tag("site_name", s.SiteName),
			tag("source", s.SourceName),
		}

		_ = r.reportCount(metricName("tx_packets"), dpi.TxPackets.Int64(), tags)
		_ = r.reportCount(metricName("rx_packets"), dpi.RxPackets.Int64(), tags)
		_ = r.reportCount(metricName("tx_bytes"), dpi.TxBytes.Int64(), tags)
		_ = r.reportCount(metricName("rx_bytes"), dpi.RxBytes.Int64(), tags)
	}
}
