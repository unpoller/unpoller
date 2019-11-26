package promunifi

import (
	"golift.io/unifi"
)

type udm struct {
}

func descUDM(ns string) *udm {
	return &udm{}
}

func (u *unifiCollector) exportUDMs(udms []*unifi.UDM, r *Report) {
	//	for _, d := range udms {
}
