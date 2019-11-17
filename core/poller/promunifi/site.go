package promunifi

import (
	"golift.io/unifi"
)

type site struct {
}

func descSite(ns string) *site {
	return &site{}
}

// exportSite exports Network Site Data
func (u *unifiCollector) exportSite(s *unifi.Site) []*metricExports {
	return nil
}
