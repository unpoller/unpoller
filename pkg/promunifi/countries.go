package promunifi

import (
	"github.com/flaticols/countrycodes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unpoller/unifi/v5"
)

type ucountrytraffic struct {
	RxBytes *prometheus.Desc
	TxBytes *prometheus.Desc
}

func descCountryTraffic(ns string) *ucountrytraffic {
	labels := []string{
		"code",
		"name",
	}

	return &ucountrytraffic{
		RxBytes: prometheus.NewDesc(ns+"receive_bytes_total", "Country Receive Bytes", labels, nil),
		TxBytes: prometheus.NewDesc(ns+"transmit_bytes_total", "Country Transmit Bytes", labels, nil),
	}
}

func (u *promUnifi) exportCountryTraffic(r report, v any) {
	s, ok := v.(*unifi.UsageByCountry)
	if !ok {
		u.LogErrorf("invalid type given to CountryTraffic: %T", v)
		return
	}
	name, ok := countrycodes.Alpha2ToName(s.Country)
	if !ok {
		name = "Unknown"
	}
	labels := []string{s.Country, name}
	r.send([]*metric{
		{u.CountryTraffic.RxBytes, counter, s.BytesReceived, labels},
		{u.CountryTraffic.TxBytes, counter, s.BytesTransmitted, labels},
	})
}
