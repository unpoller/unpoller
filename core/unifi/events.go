package unifi

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var (
	errNoSiteProvided   = fmt.Errorf("site must not be nil or empty")
	errInvalidTimeRange = fmt.Errorf("only 0, 1 or 2 times may be provided to timeRange")
)

const (
	eventLimit = 50000
)

// GetEvents returns a response full of UniFi Events from multiple sites.
// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
// Events between start and end are returned. End defaults to time.Now().
func (u *Unifi) GetEvents(sites []*Site, timeRange ...time.Time) ([]*Event, error) {
	data := make([]*Event, 0)

	for _, site := range sites {
		response, err := u.GetSiteEvents(site, timeRange...)
		if err != nil {
			return data, err
		}

		data = append(data, response...)
	}

	return data, nil
}

// GetSiteEvents retreives the events from a single site.
// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
// Events between start and end are returned. End defaults to time.Now().
func (u *Unifi) GetSiteEvents(site *Site, timeRange ...time.Time) ([]*Event, error) { // nolint: dupl
	if site == nil || site.Name == "" {
		return nil, errNoSiteProvided
	}

	u.DebugLog("Polling Controller, retreiving UniFi Events, site %s (%s)", site.Name, site.Desc)

	var (
		path  = fmt.Sprintf(APIEventPath, site.Name)
		event struct {
			Data []*Event `json:"data"`
		}
	)

	if params, err := makeEventParams(timeRange...); err != nil {
		return event.Data, err
	} else if err = u.GetData(path, &event, params); err != nil {
		return event.Data, err
	}

	for i := range event.Data {
		// Add special SourceName value.
		event.Data[i].SourceName = u.URL
		// Add the special "Site Name" to each event. This becomes a Grafana filter somewhere.
		event.Data[i].SiteName = site.Desc + " (" + site.Name + ")"
	}

	return event.Data, nil
}

func makeEventParams(timeRange ...time.Time) (string, error) {
	type eventReq struct {
		Start  int64  `json:"start,omitempty"`
		End    int64  `json:"end,omitempty"`
		Limit  int    `json:"_limit,omitempty"`
		Within int    `json:"within"`
		Sort   string `json:"_sort"`
	}

	rp := eventReq{Limit: eventLimit, Sort: "-time", Within: 1}

	switch len(timeRange) {
	case 0:
		rp.End = time.Now().Unix() * int64(time.Microsecond)
	case 1:
		rp.Start = timeRange[0].Unix() * int64(time.Microsecond)
		rp.End = time.Now().Unix() * int64(time.Microsecond)
	case 2: // nolint: gomnd
		rp.Start = timeRange[0].Unix() * int64(time.Microsecond)
		rp.End = timeRange[1].Unix() * int64(time.Microsecond)
	default:
		return "", errInvalidTimeRange
	}

	params, err := json.Marshal(&rp)

	log.Println(string(params))

	return string(params), err
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
	GeoIP
}

// GeoIP is a struct in a struct to deal with weird UniFi output.
type GeoIP struct {
	Asn           int64   `json:"asn"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	City          string  `json:"city"`
	ContinentCode string  `json:"continent_code"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	Organization  string  `json:"organization"`
}

// UnmarshalJSON is required because sometimes the unifi api returns
// an empty array instead of a struct filled with data.
func (v *IPGeo) UnmarshalJSON(data []byte) error {
	if string(data) == "[]" {
		return nil // it's empty
	}

	return json.Unmarshal(data, &v.GeoIP)
}
