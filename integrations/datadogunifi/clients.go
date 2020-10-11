package datadogunifi

import (
	"github.com/unifi-poller/unifi"
)

// reportClient generates Unifi Client datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *DatadogUnifi) reportClient(r report, s *unifi.Client) { // nolint: funlen
	tags := []string{
		tag("mac", s.Mac),
		tag("site_name", s.SiteName),
		tag("source", s.SourceName),
		tag("ap_name", s.ApName),
		tag("gw_name", s.GwName),
		tag("sw_name", s.SwName),
		tag("oui", s.Oui),
		tag("radio_name", s.RadioName),
		tag("radio", s.Radio),
		tag("radio_proto", s.RadioProto),
		tag("name", s.Name),
		tag("fixed_ip", s.FixedIP),
		tag("sw_port", s.SwPort.Txt),
		tag("os_class", s.OsClass.Txt),
		tag("os_name", s.OsName.Txt),
		tag("dev_cat", s.DevCat.Txt),
		tag("dev_id", s.DevID.Txt),
		tag("dev_vendor", s.DevVendor.Txt),
		tag("dev_family", s.DevFamily.Txt),
		tag("is_wired", s.IsWired.Txt),
		tag("is_guest", s.IsGuest.Txt),
		tag("use_fixedip", s.UseFixedIP.Txt),
		tag("channel", s.Channel.Txt),
		tag("vlan", s.Vlan.Txt),
		tag("hostname", s.Name),
		tag("radio_desc", s.RadioDescription),
		tag("ip", s.IP),
		tag("essid", s.Essid),
		tag("bssid", s.Bssid),
	}

	data := map[string]float64{
		"anomalies":        float64(s.Anomalies),
		"channel":          s.Channel.Val,
		"satisfaction":     s.Satisfaction.Val,
		"bytes_r":          float64(s.BytesR),
		"ccq":              float64(s.Ccq),
		"noise":            float64(s.Noise),
		"roam_count":       float64(s.RoamCount),
		"rssi":             float64(s.Rssi),
		"rx_bytes":         float64(s.RxBytes),
		"rx_bytes_r":       float64(s.RxBytesR),
		"rx_packets":       float64(s.RxPackets),
		"rx_rate":          float64(s.RxRate),
		"signal":           float64(s.Signal),
		"tx_bytes":         float64(s.TxBytes),
		"tx_bytes_r":       float64(s.TxBytesR),
		"tx_packets":       float64(s.TxPackets),
		"tx_retries":       float64(s.TxRetries),
		"tx_power":         float64(s.TxPower),
		"tx_rate":          float64(s.TxRate),
		"uptime":           float64(s.Uptime),
		"wifi_tx_attempts": float64(s.WifiTxAttempts),
		"wired-rx_bytes":   float64(s.WiredRxBytes),
		"wired-rx_bytes-r": float64(s.WiredRxBytesR),
		"wired-rx_packets": float64(s.WiredRxPackets),
		"wired-tx_bytes":   float64(s.WiredTxBytes),
		"wired-tx_bytes-r": float64(s.WiredTxBytesR),
		"wired-tx_packets": float64(s.WiredTxPackets),
	}
	metricName := metricNamespace("clients")
	reportGaugeForMap(r, metricName, data, tags)
}

// totalsDPImap: controller, site, name (app/cat name), dpi.
type totalsDPImap map[string]map[string]map[string]unifi.DPIData

func (u *DatadogUnifi) reportClientDPI(r report, s *unifi.DPITable, appTotal, catTotal totalsDPImap) {
	for _, dpi := range s.ByApp {
		category := unifi.DPICats.Get(dpi.Cat)
		application := unifi.DPIApps.GetApp(dpi.Cat, dpi.App)
		fillDPIMapTotals(appTotal, application, s.SourceName, s.SiteName, dpi)
		fillDPIMapTotals(catTotal, category, s.SourceName, s.SiteName, dpi)

		tags := []string{
			tag("category", category),
			tag("application", application),
			tag("name", s.Name),
			tag("mac", s.MAC),
			tag("site_name", s.SiteName),
			tag("source", s.SourceName),
		}
		data := map[string]float64{
			"tx_packets": float64(dpi.TxPackets),
			"rx_packets": float64(dpi.RxPackets),
			"tx_bytes":   float64(dpi.TxBytes),
			"rx_bytes":   float64(dpi.RxBytes),
		}
		metricName := metricNamespace("clientdpi")
		reportGaugeForMap(r, metricName, data, tags)
	}
}

// fillDPIMapTotals fills in totals for categories and applications. maybe clients too.
// This allows less processing in InfluxDB to produce total transfer data per cat or app.
func fillDPIMapTotals(m totalsDPImap, name, controller, site string, dpi unifi.DPIData) {
	if m[controller] == nil {
		m[controller] = make(map[string]map[string]unifi.DPIData)
	}

	if m[controller][site] == nil {
		m[controller][site] = make(map[string]unifi.DPIData)
	}

	existing := m[controller][site][name]
	existing.TxPackets += dpi.TxPackets
	existing.RxPackets += dpi.RxPackets
	existing.TxBytes += dpi.TxBytes
	existing.RxBytes += dpi.RxBytes
	m[controller][site][name] = existing
}

func reportClientDPItotals(r report, appTotal, catTotal totalsDPImap) {
	type all []struct {
		kind string
		val  totalsDPImap
	}

	// This produces 7000+ metrics per site. Disabled for now.
	if appTotal != nil {
		appTotal = nil
	}

	// This can allow us to aggregate other data types later, like `name` or `mac`, or anything else unifi adds.
	a := all{
		// This produces 7000+ metrics per site. Disabled for now.
		{
			kind: "application",
			val:  appTotal,
		},
		{
			kind: "category",
			val:  catTotal,
		},
	}

	for _, k := range a {
		for controller, s := range k.val {
			for site, c := range s {
				for name, m := range c {
					tags := []string{
						tag("site_name", site),
						tag("source", controller),
						tag("name", name),
					}
					data := map[string]float64{
						"tx_packets": float64(m.TxPackets),
						"rx_packets": float64(m.RxPackets),
						"tx_bytes":   float64(m.TxBytes),
						"rx_bytes":   float64(m.RxBytes),
					}
					metricName := metricNamespace("clientdpi.totals")
					reportGaugeForMap(r, metricName, data, tags)
				}
			}
		}
	}
}
