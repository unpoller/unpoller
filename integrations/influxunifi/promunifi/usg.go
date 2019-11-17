package promunifi

import (
	"golift.io/unifi"
)

type usg struct {
}

func descUSG(ns string) *usg {
	return &usg{}
}

// exportUSG Exports Security Gateway Data
func (u *unifiCollector) exportUSG(s *unifi.USG) []*metricExports {
	return nil
}
