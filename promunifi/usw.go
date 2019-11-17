package promunifi

import (
	"golift.io/unifi"
)

type usw struct {
}

func descUSW(ns string) *usw {
	return &usw{}
}

// exportUSW exports Network Switch Data
func (u *unifiCollector) exportUSW(s *unifi.USW) []*metricExports {
	return nil
}
