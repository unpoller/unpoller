package unifi

import (
	"fmt"
	"sort"
	"time"
)

type Alarm struct {
	Archived              FlexBool  `json:"archived"`
	DestPort              int       `json:"dest_port"`
	SrcPort               int       `json:"src_port"`
	FlowID                int64     `json:"flow_id"`
	InnerAlertGID         int64     `json:"inner_alert_gid"`
	InnerAlertRev         int64     `json:"inner_alert_rev"`
	InnerAlertSeverity    int64     `json:"inner_alert_severity"`
	InnerAlertSignatureID int64     `json:"inner_alert_signature_id"`
	Time                  int64     `json:"time"`
	Timestamp             int64     `json:"timestamp"`
	Datetime              time.Time `json:"datetime"`
	HandledTime           time.Time `json:"handled_time,omitempty"`
	AppProto              string    `json:"app_proto,omitempty"`
	Catname               string    `json:"catname"`
	DestIP                string    `json:"dest_ip"`
	DstMAC                string    `json:"dst_mac"`
	DstIPASN              string    `json:"dstipASN,omitempty"`
	DstIPCountry          string    `json:"dstipCountry,omitempty"`
	EventType             string    `json:"event_type"`
	HandledAdminID        string    `json:"handled_admin_id,omitempty"`
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
	SrcIPASN              string    `json:"srcipASN,omitempty"`
	SrcIPCountry          string    `json:"srcipCountry,omitempty"`
	SrcMAC                string    `json:"src_mac"`
	Subsystem             string    `json:"subsystem"`
	UniqueAlertID         string    `json:"unique_alertid"`
	USGIP                 string    `json:"usgip"`
	USGIPASN              string    `json:"usgipASN"`
	USGIPCountry          string    `json:"usgipCountry"`
	TxID                  FlexInt   `json:"tx_id,omitempty"`
	DestIPGeo             IPGeo     `json:"dstipGeo"`
	SourceIPGeo           IPGeo     `json:"usgipGeo"`
	USGIPGeo              IPGeo     `json:"srcipGeo,omitempty"`
}

// GetAlarms returns Alarms for a list of Sites.
func (u *Unifi) GetAlarms(sites []*Site) ([]*Alarm, error) {
	data := []*Alarm{}

	for _, site := range sites {
		response, err := u.GetAlarmsSite(site)
		if err != nil {
			return data, err
		}

		data = append(data, response...)
	}

	return data, nil
}

// GetAlarmsSite retreives the Alarms for a single Site.
func (u *Unifi) GetAlarmsSite(site *Site) ([]*Alarm, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for Alarms, site %s", site.SiteName)

	var (
		path   = fmt.Sprintf(APIEventPathAlarms, site.Name)
		alarms struct {
			Data alarms `json:"data"`
		}
	)

	if err := u.GetData(path, &alarms, ""); err != nil {
		return alarms.Data, err
	}

	for i := range alarms.Data {
		// Add special SourceName value.
		alarms.Data[i].SourceName = u.URL
		// Add the special "Site Name" to each event. This becomes a Grafana filter somewhere.
		alarms.Data[i].SiteName = site.SiteName
	}

	sort.Sort(alarms.Data)

	return alarms.Data, nil
}

// alarms satisfies the sort.Sort Interface.
type alarms []*Alarm

// Len satisfies sort.Interface.
func (a alarms) Len() int {
	return len(a)
}

// Swap satisfies sort.Interface.
func (a alarms) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less satisfies sort.Interface. Sort our list by Datetime.
func (a alarms) Less(i, j int) bool {
	return a[i].Datetime.Before(a[j].Datetime)
}
