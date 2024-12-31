package datadogunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchClient generates Unifi Client datapoints for Datadog.
// These points can be passed directly to Datadog.
func (u *DatadogUnifi) batchClient(r report, s *unifi.Client) { // nolint: funlen
	tags := map[string]string{
		"mac":          s.Mac,
		"site_name":    s.SiteName,
		"source":       s.SourceName,
		"ap_name":      s.ApName,
		"gw_name":      s.GwName,
		"sw_name":      s.SwName,
		"oui":          s.Oui,
		"radio_name":   s.RadioName,
		"radio":        s.Radio,
		"radio_proto":  s.RadioProto,
		"name":         s.Name,
		"fixed_ip":     s.FixedIP,
		"sw_port":      s.SwPort.Txt,
		"os_class":     s.OsClass.Txt,
		"os_name":      s.OsName.Txt,
		"dev_cat":      s.DevCat.Txt,
		"dev_id":       s.DevID.Txt,
		"dev_vendor":   s.DevVendor.Txt,
		"dev_family":   s.DevFamily.Txt,
		"is_wired":     s.IsWired.Txt,
		"is_guest":     s.IsGuest.Txt,
		"use_fixed_ip": s.UseFixedIP.Txt,
		"channel":      s.Channel.Txt,
		"vlan":         s.Vlan.Txt,
		"hostname":     s.Name,
		"essid":        s.Essid,
		"bssid":        s.Bssid,
		"ip":           s.IP,
	}

	powerSaveEnabled := 0.0
	if s.PowersaveEnabled.Val {
		powerSaveEnabled = 1.0
	}

	data := map[string]float64{
		"anomalies":         s.Anomalies.Val,
		"channel":           s.Channel.Val,
		"satisfaction":      s.Satisfaction.Val,
		"bytes_r":           s.BytesR.Val,
		"ccq":               s.Ccq.Val,
		"noise":             s.Noise.Val,
		"powersave_enabled": powerSaveEnabled,
		"roam_count":        s.RoamCount.Val,
		"rssi":              s.Rssi.Val,
		"rx_bytes":          s.RxBytes.Val,
		"rx_bytes_r":        s.RxBytesR.Val,
		"rx_packets":        s.RxPackets.Val,
		"rx_rate":           s.RxRate.Val,
		"signal":            s.Signal.Val,
		"tx_bytes":          s.TxBytes.Val,
		"tx_bytes_r":        s.TxBytesR.Val,
		"tx_packets":        s.TxPackets.Val,
		"tx_retries":        s.TxRetries.Val,
		"tx_power":          s.TxPower.Val,
		"tx_rate":           s.TxRate.Val,
		"uptime":            s.Uptime.Val,
		"wifi_tx_attempts":  s.WifiTxAttempts.Val,
		"wired_rx_bytes":    s.WiredRxBytes.Val,
		"wired_rx_bytes-r":  s.WiredRxBytesR.Val,
		"wired_rx_packets":  s.WiredRxPackets.Val,
		"wired_tx_bytes":    s.WiredTxBytes.Val,
		"wired_tx_bytes-r":  s.WiredTxBytesR.Val,
		"wired_tx_packets":  s.WiredTxPackets.Val,
	}

	metricName := metricNamespace("clients")

	reportGaugeForFloat64Map(r, metricName, data, tags)
}

// totalsDPImap: controller, site, name (app/cat name), dpi.
type totalsDPImap map[string]map[string]map[string]unifi.DPIData

func (u *DatadogUnifi) batchClientDPI(r report, v any, appTotal, catTotal totalsDPImap) {
	s, ok := v.(*unifi.DPITable)
	if !ok {
		u.LogErrorf("invalid type given to batchClientDPI: %T", v)

		return
	}

	for _, dpi := range s.ByApp {
		category := unifi.DPICats.Get(dpi.Cat.Int())
		application := unifi.DPIApps.GetApp(dpi.Cat.Int(), dpi.App.Int())
		fillDPIMapTotals(appTotal, application, s.SourceName, s.SiteName, dpi)
		fillDPIMapTotals(catTotal, category, s.SourceName, s.SiteName, dpi)

		tags := map[string]string{
			"category":    category,
			"application": application,
			"name":        s.Name,
			"mac":         s.MAC,
			"site_name":   s.SiteName,
			"source":      s.SourceName,
		}

		data := map[string]float64{
			"tx_packets": dpi.TxPackets.Val,
			"rx_packets": dpi.RxPackets.Val,
			"tx_bytes":   dpi.TxBytes.Val,
			"rx_bytes":   dpi.RxBytes.Val,
		}

		metricName := metricNamespace("client_dpi")

		reportGaugeForFloat64Map(r, metricName, data, tags)
	}
}

// fillDPIMapTotals fills in totals for categories and applications. maybe clients too.
// This allows less processing in Datadog to produce total transfer data per cat or app.
func fillDPIMapTotals(m totalsDPImap, name, controller, site string, dpi unifi.DPIData) {
	if m[controller] == nil {
		m[controller] = make(map[string]map[string]unifi.DPIData)
	}

	if m[controller][site] == nil {
		m[controller][site] = make(map[string]unifi.DPIData)
	}

	existing := m[controller][site][name]
	existing.TxPackets.Add(&dpi.TxPackets)
	existing.RxPackets.Add(&dpi.RxPackets)
	existing.TxBytes.Add(&dpi.TxBytes)
	existing.RxBytes.Add(&dpi.RxBytes)
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
					tags := map[string]string{
						"category":    "TOTAL",
						"application": "TOTAL",
						"name":        "TOTAL",
						"mac":         "TOTAL",
						"site_name":   site,
						"source":      controller,
					}
					tags[k.kind] = name

					data := map[string]float64{
						"tx_packets": m.TxPackets.Val,
						"rx_packets": m.RxPackets.Val,
						"tx_bytes":   m.TxBytes.Val,
						"rx_bytes":   m.RxBytes.Val,
					}

					metricName := metricNamespace("client_dpi")

					reportGaugeForFloat64Map(r, metricName, data, tags)
				}
			}
		}
	}
}
