package promunifi

import (
	"github.com/unpoller/unifi/v5"
)

// exportUCI is a collection of stats from UCI.
func (u *promUnifi) exportUCI(r report, d *unifi.UCI) {
	if !d.Adopted.Val || d.Locating.Val {
		return
	}

	var sw *unifi.Sw
	if d.Stat != nil {
		sw = d.Stat.Sw
	}

	labels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	infoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}
	// Shared data (all devices do this).
	u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)

	if d.SysStats != nil && d.SystemStats != nil {
		u.exportSYSstats(r, labels, *d.SysStats, *d.SystemStats)
	}

	// Switch Data
	u.exportUSWstats(r, labels, sw)
	// Dream Machine System Data.
	r.send([]*metric{
		{u.Device.Info, gauge, 1.0, append(labels, infoLabels...)},
		{u.Device.Uptime, gauge, d.Uptime, labels},
	})
}
