package promunifi

import "github.com/unpoller/unifi/v5"

// exportUDB exports metrics for UDB (UniFi Device Bridge) devices.
// The UDB range includes UDB-Switch, UDB-Pro, UDB-Pro-Sector.
// UDB-Switch is a hybrid device combining switch ports (8 PoE ports)
// with WiFi 7 wireless bridge capability (5GHz + 6GHz radios).
func (u *promUnifi) exportUDB(r report, d *unifi.UDB) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}

	baseLabels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	baseInfoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}

	u.exportWithTags(r, d.Tags, func(tagLabels []string) {
		tag := tagLabels[0]
		labels := append(baseLabels, tag)
		infoLabels := append(baseInfoLabels, tag)

		// Export switch stats (reuse USW functions)
		u.exportUSWstats(r, labels, d.Stat.Sw)
		u.exportPRTtable(r, labels, d.PortTable)

		// Export wireless stats (reuse UAP functions)
		u.exportVAPtable(r, labels, d.VapTable)
		u.exportRADtable(r, labels, d.RadioTable, d.RadioTableStats)

		// Common device stats
		u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)
		u.exportSYSstats(r, labels, d.SysStats, d.SystemStats)
		u.exportSTAcount(r, labels, d.UserNumSta, d.GuestNumSta)

		r.send([]*metric{
			{u.Device.Info, gauge, 1.0, append(baseLabels, infoLabels...)},
			{u.Device.Uptime, gauge, d.Uptime, labels},
			{u.Device.Upgradeable, gauge, d.Upgradable.Val, labels},
		})

		// Temperature and fan
		if d.HasTemperature.Val {
			r.send([]*metric{{u.Device.Temperature, gauge, d.GeneralTemperature, append(labels, "general", "board")}})
		}

		if d.HasFan.Val {
			r.send([]*metric{{u.Device.FanLevel, gauge, d.FanLevel, labels}})
		}

		if d.TotalMaxPower.Txt != "" {
			r.send([]*metric{{u.Device.TotalMaxPower, gauge, d.TotalMaxPower, labels}})
		}
	})
}
