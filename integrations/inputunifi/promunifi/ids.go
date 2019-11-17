package promunifi

import (
	"golift.io/unifi"
)

/* The IDS data does not go into prometheus cleanly. This probably wont happen. */

type ids struct {
}

func descIDS(ns string) *ids {
	return &ids{}
}

// exportIDS exports Intrusion Detection System Data
func (u *unifiCollector) exportIDS(i *unifi.IDS) []*metricExports {
	return nil
}
