package datadogunifi

import (
	"github.com/unpoller/unifi"
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
		"anomalies":         float64(s.Anomalies),
		"channel":           s.Channel.Val,
		"satisfaction":      s.Satisfaction.Val,
		"bytes_r":           float64(s.BytesR),
		"ccq":               float64(s.Ccq),
		"noise":             float64(s.Noise),
		"powersave_enabled": powerSaveEnabled,
		"roam_count":        float64(s.RoamCount),
		"rssi":              float64(s.Rssi),
		"rx_bytes":          float64(s.RxBytes),
		"rx_bytes_r":        float64(s.RxBytesR),
		"rx_packets":        float64(s.RxPackets),
		"rx_rate":           float64(s.RxRate),
		"signal":            float64(s.Signal),
		"tx_bytes":          float64(s.TxBytes),
		"tx_bytes_r":        float64(s.TxBytesR),
		"tx_packets":        float64(s.TxPackets),
		"tx_retries":        float64(s.TxRetries),
		"tx_power":          float64(s.TxPower),
		"tx_rate":           float64(s.TxRate),
		"uptime":            float64(s.Uptime),
		"wifi_tx_attempts":  float64(s.WifiTxAttempts),
		"wired_rx_bytes":    float64(s.WiredRxBytes),
		"wired_rx_bytes-r":  float64(s.WiredRxBytesR),
		"wired_rx_packets":  float64(s.WiredRxPackets),
		"wired_tx_bytes":    float64(s.WiredTxBytes),
		"wired_tx_bytes-r":  float64(s.WiredTxBytesR),
		"wired_tx_packets":  float64(s.WiredTxPackets),
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
		category := unifi.DPICats.Get(dpi.Cat)
		application := unifi.DPIApps.GetApp(dpi.Cat, dpi.App)
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
			"tx_packets": float64(dpi.TxPackets),
			"rx_packets": float64(dpi.RxPackets),
			"tx_bytes":   float64(dpi.TxBytes),
			"rx_bytes":   float64(dpi.RxBytes),
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
						"tx_packets": float64(m.TxPackets),
						"rx_packets": float64(m.RxPackets),
						"tx_bytes":   float64(m.TxBytes),
						"rx_bytes":   float64(m.RxBytes),
					}

					metricName := metricNamespace("client_dpi")

					reportGaugeForFloat64Map(r, metricName, data, tags)
				}
			}
		}
	}
}
