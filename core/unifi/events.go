package unifi

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

var (
	ErrNoSiteProvided   = fmt.Errorf("site must not be nil or empty")
	ErrInvalidTimeRange = fmt.Errorf("only 0, 1 or 2 times may be provided to timeRange")
)

const (
	eventLimit = 50000
)

// GetEvents returns a response full of UniFi Events for the last 1 hour from multiple sites.
func (u *Unifi) GetEvents(sites []*Site, hours time.Duration) ([]*Event, error) {
	data := make([]*Event, 0)

	for _, site := range sites {
		response, err := u.GetSiteEvents(site, hours)
		if err != nil {
			return data, err
		}

		data = append(data, response...)
	}

	return data, nil
}

// GetSiteEvents retrieves the last 1 hour's worth of events from a single site.
func (u *Unifi) GetSiteEvents(site *Site, hours time.Duration) ([]*Event, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	if hours < time.Hour {
		hours = time.Hour
	}

	u.DebugLog("Polling Controller, retreiving UniFi Events, site %s", site.SiteName)

	var (
		path   = fmt.Sprintf(APIEventPath, site.Name)
		params = fmt.Sprintf(`{"_limit":%d,"within":%d,"_sort":"-time"}}`,
			eventLimit, int(hours.Round(time.Hour).Hours()))
		event struct {
			Data events `json:"data"`
		}
	)

	if err := u.GetData(path, &event, params); err != nil {
		return event.Data, err
	}

	for i := range event.Data {
		// Add special SourceName value.
		event.Data[i].SourceName = u.URL
		// Add the special "Site Name" to each event. This becomes a Grafana filter somewhere.
		event.Data[i].SiteName = site.SiteName
	}

	sort.Sort(event.Data)

	return event.Data, nil
}

// Event describes a UniFi Event.
// API Path: /api/s/default/stat/event.
type Event struct {
	IsAdmin               FlexBool  `json:"is_admin"`
	DestPort              int       `json:"dest_port"`
	SrcPort               int       `json:"src_port"`
	Bytes                 FlexInt   `json:"bytes"`
	Duration              FlexInt   `json:"duration"`
	FlowID                FlexInt   `json:"flow_id"`
	InnerAlertGID         FlexInt   `json:"inner_alert_gid"`
	InnerAlertRev         FlexInt   `json:"inner_alert_rev"`
	InnerAlertSeverity    FlexInt   `json:"inner_alert_severity"`
	InnerAlertSignatureID FlexInt   `json:"inner_alert_signature_id"`
	Channel               FlexInt   `json:"channel"`
	ChannelFrom           FlexInt   `json:"channel_from"`
	ChannelTo             FlexInt   `json:"channel_to"`
	Time                  int64     `json:"time"`
	Timestamp             int64     `json:"timestamp"`
	Datetime              time.Time `json:"datetime"`
	Admin                 string    `json:"admin"`
	Ap                    string    `json:"ap"`
	ApFrom                string    `json:"ap_from"`
	ApName                string    `json:"ap_name"`
	ApTo                  string    `json:"ap_to"`
	AppProto              string    `json:"app_proto"`
	Catname               string    `json:"catname"`
	DestIP                string    `json:"dest_ip"`
	DstMAC                string    `json:"dst_mac"`
	EventType             string    `json:"event_type"`
	Guest                 string    `json:"guest"`
	Gw                    string    `json:"gw"`
	GwName                string    `json:"gw_name"`
	Host                  string    `json:"host"`
	Hostname              string    `json:"hostname"`
	ID                    string    `json:"_id"`
	IP                    string    `json:"ip"`
	InIface               string    `json:"in_iface"`
	InnerAlertAction      string    `json:"inner_alert_action"`
	InnerAlertCategory    string    `json:"inner_alert_category"`
	InnerAlertSignature   string    `json:"inner_alert_signature"`
	Key                   string    `json:"key"`
	Msg                   string    `json:"msg"`
	Network               string    `json:"network"`
	Proto                 string    `json:"proto"`
	Radio                 string    `json:"radio"`
	RadioFrom             string    `json:"radio_from"`
	RadioTo               string    `json:"radio_to"`
	SiteID                string    `json:"site_id"`
	SiteName              string    `json:"-"`
	SourceName            string    `json:"-"`
	SrcIP                 string    `json:"src_ip"`
	SrcMAC                string    `json:"src_mac"`
	SrcIPASN              string    `json:"srcipASN"`
	SrcIPCountry          string    `json:"srcipCountry"`
	SSID                  string    `json:"ssid"`
	Subsystem             string    `json:"subsystem"`
	Sw                    string    `json:"sw"`
	SwName                string    `json:"sw_name"`
	UniqueAlertID         string    `json:"unique_alertid"`
	User                  string    `json:"user"`
	USGIP                 string    `json:"usgip"`
	USGIPASN              string    `json:"usgipASN"`
	USGIPCountry          string    `json:"usgipCountry"`
	DestIPGeo             IPGeo     `json:"dstipGeo"`
	SourceIPGeo           IPGeo     `json:"srcipGeo"`
	USGIPGeo              IPGeo     `json:"usgipGeo"`
}

// IPGeo is part of the UniFi Event data. Each event may have up to three of these.
// One for source, one for dest and one for the USG location.
type IPGeo struct {
	Asn           int64   `json:"asn"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	City          string  `json:"city"`
	ContinentCode string  `json:"continent_code"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	Organization  string  `json:"organization"`
}

// Events satisfied the sort.Interface.
type events []*Event

// Len satisfies sort.Interface.
func (e events) Len() int {
	return len(e)
}

// Swap satisfies sort.Interface.
func (e events) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// Less satisfies sort.Interface. Sort our list by date/time.
func (e events) Less(i, j int) bool {
	return e[i].Datetime.Before(e[j].Datetime)
}

// UnmarshalJSON is required because sometimes the unifi api returns
// an empty array instead of a struct filled with data.
func (v *IPGeo) UnmarshalJSON(data []byte) error {
	if string(data) == "[]" {
		return nil // it's empty
	}

	g := struct {
		Asn           int64   `json:"asn"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		City          string  `json:"city"`
		ContinentCode string  `json:"continent_code"`
		CountryCode   string  `json:"country_code"`
		CountryName   string  `json:"country_name"`
		Organization  string  `json:"organization"`
	}{}

	err := json.Unmarshal(data, &g)
	v.Asn = g.Asn
	v.Latitude = g.Latitude
	v.Longitude = g.Longitude
	v.City = g.City
	v.ContinentCode = g.ContinentCode
	v.CountryCode = g.CountryCode
	v.CountryName = g.CountryName
	v.Organization = g.Organization

	if err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	return nil
}
