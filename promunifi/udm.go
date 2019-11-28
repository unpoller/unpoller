package promunifi

import "golift.io/unifi"

type udm struct {
}

func descUDM(ns string) *udm {
	return &udm{}
}

func (u *unifiCollector) exportUDMs(r *Report) {
	if r.Metrics == nil || r.Metrics.Devices == nil || len(r.Metrics.Devices.UDMs) < 1 {
		return
	}
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for _, d := range r.Metrics.Devices.UDMs {
			u.exportUDM(r, d)
		}
	}()
}

func (u *unifiCollector) exportUDM(r *Report, d *unifi.UDM) {
	//	for _, d := range r.Metrics.Devices.UDMs {
}
