package promunifi

import (
	"golift.io/unifi"
)

type udm struct {
}

func descUDM(ns string) *udm {
	return &udm{}
}

func (u *unifiCollector) exportUDMs(udms []*unifi.UDM) (e []*metricExports) {
	for _, d := range udms {
		e = append(e, u.exportUDM(d)...)
	}
	return
}

// exportUDM exports UniFi Dream Machine (and Pro) Data
func (u *unifiCollector) exportUDM(d *unifi.UDM) []*metricExports {
	return nil
}
