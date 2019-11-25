package promunifi

import (
	"golift.io/unifi"
)

type udm struct {
}

func descUDM(ns string) *udm {
	return &udm{}
}

func (u *unifiCollector) exportUDMs(udms []*unifi.UDM, ch chan []*metricExports) {
	for _, d := range udms {
		ch <- u.exportUDM(d)
	}
}

// exportUDM exports UniFi Dream Machine (and Pro) Data
func (u *unifiCollector) exportUDM(d *unifi.UDM) []*metricExports {
	return nil
}
