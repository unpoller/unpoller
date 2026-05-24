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

	baseLabels := []string{d.Type, d.SiteName, d.Name, d.SourceName}
	baseInfoLabels := []string{d.Version, d.Model, d.Serial, d.Mac, d.IP, d.ID}
	
	u.exportWithTags(r, d.Tags, func(tagLabels []string) {
		tag := tagLabels[0]
		labels := append(baseLabels, tag)
		infoLabels := append(baseInfoLabels, tag)
		
		// Shared data (all devices do this).
		u.exportBYTstats(r, labels, d.TxBytes, d.RxBytes)

		if d.SysStats != nil && d.SystemStats != nil {
			u.exportSYSstats(r, labels, *d.SysStats, *d.SystemStats)
		}

		// Switch Data
		u.exportUSWstats(r, labels, sw)
		// Dream Machine System Data.
		r.send([]*metric{
			{u.Device.Info, gauge, 1.0, append(baseLabels, infoLabels...)},
			{u.Device.Uptime, gauge, d.Uptime, labels},
		})

		// DOCSIS / Cable Internet state.
		//
		// The controller exposes a top-level `internet` boolean on UCI devices,
		// but it is not a reliable WAN-reachability signal — it stays false
		// even when the cable link is fully up and traffic is flowing. The
		// UCI is a cable bridge with no independent internet-reachability
		// check; real WAN health lives on the upstream gateway (e.g. UDM
		// `wan1.up`).
		//
		// Instead, derive an operational gauge from the DOCSIS CI state,
		// which IS reported reliably by the controller. The full state
		// string is also exposed as a label on `*_ci_state_info`.
		if d.CiStateTable != nil {
			r.send([]*metric{
				{
					u.Device.CiStateOperational, gauge,
					d.CiStateTable.CIState == "Operational",
					labels,
				},
				{
					u.Device.CiStateInfo, gauge, 1.0,
					append(labels, d.CiStateTable.CIState, d.CiStateTable.CISwDlStatus,
						d.CiStateTable.CIMac, d.CiStateTable.CIVersion, d.CiStateTable.CIMode),
				},
			})
		}
	})
}
