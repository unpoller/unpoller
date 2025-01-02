package promunifi

import (
	"github.com/unpoller/unifi/v5"
)

// exportUBB is a collection of stats from UBB.
func (u *promUnifi) exportUBB(r report, d *unifi.UBB) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}

	//var sw *unifi.Bb
	//if d.Stat != nil {
	//	sw = d.Stat.Bb
	//}
	// unsure of what to do with this yet.

	labels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}
	// Shared data (all devices do this).
	u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)

	if d.SysStats != nil && d.SystemStats != nil {
		u.exportSYSstats(r, labels, *d.SysStats, *d.SystemStats)
	}

	// Dream Machine System Data.
	r.send([]*metric{
		{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
		{u.Device.Uptime, gauge, d.Uptime, labels},
	})

	// temperature
	r.send([]*metric{{u.Device.Temperature, gauge, d.GeneralTemperature.Val, append(labels, d.Name, "general")}})
}
