package promunifi

import (
	"golift.io/unifi"
)

type udm struct {
}

func descUDM(ns string) *udm {
	return &udm{}
}

// exportUDM exports UniFi Dream Machine (and Pro) Data
func (u *unifiCollector) exportUDM(d *unifi.UDM) []*metricExports {
	return nil
}
