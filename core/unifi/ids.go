package unifi

import (
	"fmt"
	"time"
)

// IDS holds an Intrusion Prevention System Event.
type IDS struct {
	Archived              FlexBool  `json:"archived"`
	DstIPCountry          FlexBool  `json:"dstipCountry"`
	DestPort              int       `json:"dest_port,omitempty"`
	SrcPort               int       `json:"src_port,omitempty"`
	InnerAlertRev         int64     `json:"inner_alert_rev"`
	InnerAlertSeverity    int64     `json:"inner_alert_severity"`
	InnerAlertGID         int64     `json:"inner_alert_gid"`
	InnerAlertSignatureID int64     `json:"inner_alert_signature_id"`
	FlowID                int64     `json:"flow_id"`
	Time                  int64     `json:"time"`
	Timestamp             int64     `json:"timestamp"`
	Datetime              time.Time `json:"datetime"`
	AppProto              string    `json:"app_proto,omitempty"`
	Catname               string    `json:"catname"`
	DestIP                string    `json:"dest_ip"`
	DstMAC                string    `json:"dst_mac"`
	DstIPASN              string    `json:"dstipASN"`
	EventType             string    `json:"event_type"`
	Host                  string    `json:"host"`
	ID                    string    `json:"_id"`
	InIface               string    `json:"in_iface"`
	InnerAlertAction      string    `json:"inner_alert_action"`
	InnerAlertCategory    string    `json:"inner_alert_category"`
	InnerAlertSignature   string    `json:"inner_alert_signature"`
	Key                   string    `json:"key"`
	Msg                   string    `json:"msg"`
	Proto                 string    `json:"proto"`
	SiteID                string    `json:"site_id"`
	SiteName              string    `json:"-"`
	SourceName            string    `json:"-"`
	SrcIP                 string    `json:"src_ip"`
	SrcMAC                string    `json:"src_mac"`
	SrcIPASN              string    `json:"srcipASN"`
	SrcIPCountry          string    `json:"srcipCountry"`
	Subsystem             string    `json:"subsystem"`
	UniqueAlertID         string    `json:"unique_alertid"`
	USGIP                 string    `json:"usgip"`
	USGIPASN              string    `json:"usgipASN"`
	USGIPCountry          string    `json:"usgipCountry"`
	DestIPGeo             IPGeo     `json:"dstipGeo"`
	SourceIPGeo           IPGeo     `json:"srcipGeo"`
	USGIPGeo              IPGeo     `json:"usgipGeo"`
}

// GetIDS returns Intrusion Detection Systems events for a list of Sites.
// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
// Events between start and end are returned. End defaults to time.Now().
func (u *Unifi) GetIDS(sites []*Site, timeRange ...time.Time) ([]*IDS, error) {
	data := []*IDS{}

	for _, site := range sites {
		response, err := u.GetIDSSite(site, timeRange...)
		if err != nil {
			return data, err
		}

		data = append(data, response...)
	}

	return data, nil
}

// GetIDSSite retreives the Intrusion Detection System Data for a single Site.
// timeRange may have a length of 0, 1 or 2. The first time is Start, the second is End.
// Events between start and end are returned. End defaults to time.Now().
func (u *Unifi) GetIDSSite(site *Site, timeRange ...time.Time) ([]*IDS, error) { // nolint: dupl
	if site == nil || site.Name == "" {
		return nil, errNoSiteProvided
	}

	u.DebugLog("Polling Controller for IDS Events, site %s (%s)", site.Name, site.Desc)

	var (
		path = fmt.Sprintf(APIEventPathIDS, site.Name)
		ids  struct {
			Data []*IDS `json:"data"`
		}
	)

	if params, err := makeEventParams(timeRange...); err != nil {
		return ids.Data, err
	} else if err = u.GetData(path, &ids, params); err != nil {
		return ids.Data, err
	}

	for i := range ids.Data {
		// Add special SourceName value.
		ids.Data[i].SourceName = u.URL
		// Add the special "Site Name" to each event. This becomes a Grafana filter somewhere.
		ids.Data[i].SiteName = site.Desc + " (" + site.Name + ")"
	}

	return ids.Data, nil
}
