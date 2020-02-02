package influxunifi

import (
	"github.com/unifi-poller/unifi"
)

// batchClient generates Unifi Client datapoints for InfluxDB.
// These points can be passed directly to influx.
func (u *InfluxUnifi) batchClient(r report, s *unifi.Client) {
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
	fields := map[string]interface{}{
		"anomalies":        s.Anomalies,
		"ip":               s.IP,
		"essid":            s.Essid,
		"bssid":            s.Bssid,
		"channel":          s.Channel.Val,
		"hostname":         s.Name,
		"radio_desc":       s.RadioDescription,
		"satisfaction":     s.Satisfaction.Val,
		"bytes_r":          s.BytesR,
		"ccq":              s.Ccq,
		"noise":            s.Noise,
		"note":             s.Note,
		"roam_count":       s.RoamCount,
		"rssi":             s.Rssi,
		"rx_bytes":         s.RxBytes,
		"rx_bytes_r":       s.RxBytesR,
		"rx_packets":       s.RxPackets,
		"rx_rate":          s.RxRate,
		"signal":           s.Signal,
		"tx_bytes":         s.TxBytes,
		"tx_bytes_r":       s.TxBytesR,
		"tx_packets":       s.TxPackets,
		"tx_retries":       s.TxRetries,
		"tx_power":         s.TxPower,
		"tx_rate":          s.TxRate,
		"uptime":           s.Uptime,
		"wifi_tx_attempts": s.WifiTxAttempts,
		"wired-rx_bytes":   s.WiredRxBytes,
		"wired-rx_bytes-r": s.WiredRxBytesR,
		"wired-rx_packets": s.WiredRxPackets,
		"wired-tx_bytes":   s.WiredTxBytes,
		"wired-tx_bytes-r": s.WiredTxBytesR,
		"wired-tx_packets": s.WiredTxPackets,
	}

	r.send(&metric{Table: "clients", Tags: tags, Fields: fields})
}

// totalsDPImap: controller, site, name (app/cat name), dpi
type totalsDPImap map[string]map[string]map[string]unifi.DPIData

func (u *InfluxUnifi) batchClientDPI(r report, s *unifi.DPITable, appTotal, catTotal totalsDPImap) {
	for _, dpi := range s.ByApp {
		category := unifi.DPICats.Get(dpi.Cat)
		application := unifi.DPIApps.GetApp(dpi.Cat, dpi.App)
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
			Fields: map[string]interface{}{
				"tx_packets": dpi.TxPackets,
				"rx_packets": dpi.RxPackets,
				"tx_bytes":   dpi.TxBytes,
				"rx_bytes":   dpi.RxBytes,
			}},
		)
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
						Fields: map[string]interface{}{
							"tx_packets": m.TxPackets,
							"rx_packets": m.RxPackets,
							"tx_bytes":   m.TxBytes,
							"rx_bytes":   m.RxBytes,
						},
					}
					newMetric.Tags[k.kind] = name

					r.send(newMetric)
				}
			}
		}
	}
}
