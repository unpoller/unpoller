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
		"region",
		"sub_region",
		"site_name",
		"source",
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
	country, ok := countrycodes.GetByAlpha2(s.Country)
	name := "Unknown"
	region := "Unknown"
	subRegion := "Unknown"
	if ok {
		name = country.Name
		region = country.Region
		subRegion = country.SubRegion
	}
	if s.Country == "GB" || s.Country == "UK" {
		name = "United Kingdom" // Because the name is so long otherwise
	}
	labels := []string{s.Country, name, region, subRegion, s.TrafficSite.SiteName, s.TrafficSite.SourceName}
	r.send([]*metric{
		{u.CountryTraffic.RxBytes, counter, s.BytesReceived, labels},
		{u.CountryTraffic.TxBytes, counter, s.BytesTransmitted, labels},
	})
}
