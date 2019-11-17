package promunifi

import (
	"golift.io/unifi"
)

type uap struct {
}

func descUAP(ns string) *uap {
	return &uap{}
}

// exportUAP exports Access Point Data
func (u *unifiCollector) exportUAP(a *unifi.UAP) []*metricExports {
	return nil
}
