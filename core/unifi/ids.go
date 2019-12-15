package unifi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// IDSList contains a list that contains all of the IDS Events on a controller.
type IDSList []*IDS

// IDS holds an Intrusion Prevention System Event.
type IDS struct {
	SourceName    string   `json:"-"`
	ID            string   `json:"_id"`
	Archived      FlexBool `json:"archived"`
	Timestamp     int64    `json:"timestamp"`
	FlowID        int64    `json:"flow_id"`
	InIface       string   `json:"in_iface"`
	EventType     string   `json:"event_type"`
	SrcIP         string   `json:"src_ip"`
	SrcMac        string   `json:"src_mac"`
	SrcPort       int      `json:"src_port,omitempty"`
	DestIP        string   `json:"dest_ip"`
	DstMac        string   `json:"dst_mac"`
	DestPort      int      `json:"dest_port,omitempty"`
	Proto         string   `json:"proto"`
	AppProto      string   `json:"app_proto,omitempty"`
	Host          string   `json:"host"`
	Usgip         string   `json:"usgip"`
	UniqueAlertid string   `json:"unique_alertid"`
	SrcipCountry  string   `json:"srcipCountry"`
	DstipCountry  FlexBool `json:"dstipCountry"`
	UsgipCountry  string   `json:"usgipCountry"`
	SrcipGeo      struct {
		ContinentCode string  `json:"continent_code"`
		CountryCode   string  `json:"country_code"`
		CountryCode3  string  `json:"country_code3"`
		CountryName   string  `json:"country_name"`
		Region        string  `json:"region"`
		City          string  `json:"city"`
		PostalCode    string  `json:"postal_code"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		DmaCode       int64   `json:"dma_code"`
		AreaCode      int64   `json:"area_code"`
	} `json:"srcipGeo"`
	DstipGeo bool `json:"dstipGeo"`
	UsgipGeo struct {
		ContinentCode string  `json:"continent_code"`
		CountryCode   string  `json:"country_code"`
		CountryCode3  string  `json:"country_code3"`
		CountryName   string  `json:"country_name"`
		Region        string  `json:"region"`
		City          string  `json:"city"`
		PostalCode    string  `json:"postal_code"`
		Latitude      float64 `json:"latitude"`
		Longitude     float64 `json:"longitude"`
		DmaCode       int64   `json:"dma_code"`
		AreaCode      int64   `json:"area_code"`
	} `json:"usgipGeo"`
	SrcipASN              string    `json:"srcipASN"`
	DstipASN              string    `json:"dstipASN"`
	UsgipASN              string    `json:"usgipASN"`
	Catname               string    `json:"catname"`
	InnerAlertAction      string    `json:"inner_alert_action"`
	InnerAlertGID         int64     `json:"inner_alert_gid"`
	InnerAlertSignatureID int64     `json:"inner_alert_signature_id"`
	InnerAlertRev         int64     `json:"inner_alert_rev"`
	InnerAlertSignature   string    `json:"inner_alert_signature"`
	InnerAlertCategory    string    `json:"inner_alert_category"`
	InnerAlertSeverity    int64     `json:"inner_alert_severity"`
	Key                   string    `json:"key"`
	Subsystem             string    `json:"subsystem"`
	SiteID                string    `json:"site_id"`
	SiteName              string    `json:"-"`
	Time                  int64     `json:"time"`
	Datetime              time.Time `json:"datetime"`
	Msg                   string    `json:"msg"`
	IcmpType              int64     `json:"icmp_type,omitempty"`
	IcmpCode              int64     `json:"icmp_code,omitempty"`
}

// GetIDS returns Intrusion Detection Systems events.
// Returns all events that happened in site between from and to.
func (u *Unifi) GetIDS(sites Sites, from, to time.Time) ([]*IDS, error) {
	data := []*IDS{}

	for _, site := range sites {
		u.DebugLog("Polling Controller for IDS/IPS Data, site %s (%s) ", site.Name, site.Desc)

		ids, err := u.GetSiteIDS(site, from, to)
		if err != nil {
			return data, err
		}
		for i := range ids {
			ids[i].SourceName = u.URL
		}
		data = append(data, ids...)
	}

	return data, nil
}

// GetSiteIDS is a helper to offload the for-loop work.
// This method retreives the Intrusion Detection System Data for 1 Site.
func (u *Unifi) GetSiteIDS(site *Site, from, to time.Time) ([]*IDS, error) {
	var response struct {
		Data []*IDS `json:"data"`
	}

	URIpath := fmt.Sprintf(APIIPSEvents, site.Name)

	params := fmt.Sprintf(`{"start":"%v000","end":"%v000","_limit":50000}`, from.Unix(), to.Unix())

	req, err := u.UniReq(URIpath, params)
	if err != nil {
		return nil, err
	}

	resp, err := u.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code from server %s", resp.Status)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	for i := range response.Data {
		response.Data[i].SiteName = site.SiteName
	}

	return response.Data, nil
}
