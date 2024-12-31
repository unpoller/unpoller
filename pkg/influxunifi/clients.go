package influxunifi

import (
	"github.com/unpoller/unifi/v5"
)

// batchClient generates Unifi Client datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchClient(r report, s *unifi.Client) { // nolint: funlen
	tags := map[string]string{
		"mac":         s.Mac,
		"site_name":   s.SiteName,
		"source":      s.SourceName,
		"ap_name":     s.ApName,
		"gw_name":     s.GwName,
		"sw_name":     s.SwName,
		"oui":         s.Oui,
		"radio_name":  s.RadioName,
		"radio":       s.Radio,
		"radio_proto": s.RadioProto,
		"name":        s.Name,
		"fixed_ip":    s.FixedIP,
		"sw_port":     s.SwPort.Txt,
		"os_class":    s.OsClass.Txt,
		"os_name":     s.OsName.Txt,
		"dev_cat":     s.DevCat.Txt,
		"dev_id":      s.DevID.Txt,
		"dev_vendor":  s.DevVendor.Txt,
		"dev_family":  s.DevFamily.Txt,
		"is_wired":    s.IsWired.Txt,
		"is_guest":    s.IsGuest.Txt,
		"use_fixedip": s.UseFixedIP.Txt,
		"channel":     s.Channel.Txt,
		"vlan":        s.Vlan.Txt,
	}

	fields := map[string]any{
		"anomalies":         s.Anomalies.Int64(),
		"ip":                s.IP,
		"essid":             s.Essid,
		"bssid":             s.Bssid,
		"channel":           s.Channel.Val,
		"hostname":          s.Name,
		"radio_desc":        s.RadioDescription,
		"satisfaction":      s.Satisfaction.Val,
		"bytes_r":           s.BytesR.Int64(),
		"ccq":               s.Ccq.Int64(),
		"noise":             s.Noise.Int64(),
		"note":              s.Note,
		"powersave_enabled": s.PowersaveEnabled,
		"roam_count":        s.RoamCount.Int64(),
		"rssi":              s.Rssi.Int64(),
		"rx_bytes":          s.RxBytes.Int64(),
		"rx_bytes_r":        s.RxBytesR.Int64(),
		"rx_packets":        s.RxPackets.Int64(),
		"rx_rate":           s.RxRate.Int64(),
		"signal":            s.Signal.Int64(),
		"tx_bytes":          s.TxBytes.Int64(),
		"tx_bytes_r":        s.TxBytesR.Int64(),
		"tx_packets":        s.TxPackets.Int64(),
		"tx_retries":        s.TxRetries.Int64(),
		"tx_power":          s.TxPower.Int64(),
		"tx_rate":           s.TxRate.Int64(),
		"uptime":            s.Uptime.Int64(),
		"wifi_tx_attempts":  s.WifiTxAttempts.Int64(),
		"wired-rx_bytes":    s.WiredRxBytes.Int64(),
		"wired-rx_bytes-r":  s.WiredRxBytesR.Int64(),
		"wired-rx_packets":  s.WiredRxPackets.Int64(),
		"wired-tx_bytes":    s.WiredTxBytes.Int64(),
		"wired-tx_bytes-r":  s.WiredTxBytesR.Int64(),
		"wired-tx_packets":  s.WiredTxPackets.Int64(),
	}

	r.send(&metric{Table: "clients", Tags: tags, Fields: fields})
}

// totalsDPImap: controller, site, name (app/cat name), dpi.
type totalsDPImap map[string]map[string]map[string]unifi.DPIData

func (u *InfluxUnifi) batchClientDPI(r report, v any, appTotal, catTotal totalsDPImap) {
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

		r.send(&metric{
			Table: "clientdpi",
			Tags: map[string]string{
				"category":    category,
				"application": application,
				"name":        s.Name,
				"mac":         s.MAC,
				"site_name":   s.SiteName,
				"source":      s.SourceName,
			},
			Fields: map[string]any{
				"tx_packets": dpi.TxPackets.Int64(),
				"rx_packets": dpi.RxPackets.Int64(),
				"tx_bytes":   dpi.TxBytes.Int64(),
				"rx_bytes":   dpi.RxBytes.Int64(),
			},
		})
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
					newMetric := &metric{
						Table: "clientdpi",
						Tags: map[string]string{
							"category":    "TOTAL",
							"application": "TOTAL",
							"name":        "TOTAL",
							"mac":         "TOTAL",
							"site_name":   site,
							"source":      controller,
						},
						Fields: map[string]any{
							"tx_packets": m.TxPackets.Int64(),
							"rx_packets": m.RxPackets.Int64(),
							"tx_bytes":   m.TxBytes.Int64(),
							"rx_bytes":   m.RxBytes.Int64(),
						},
					}
					newMetric.Tags[k.kind] = name

					r.send(newMetric)
				}
			}
		}
	}
}
