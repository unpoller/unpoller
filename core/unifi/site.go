package unifi

import (
	"fmt"
	"strings"
)

var ErrDPIDataBug = fmt.Errorf("dpi data table contains more than 1 item; please open a bug report")

// GetSites returns a list of configured sites on the UniFi controller.
func (u *Unifi) GetSites() ([]*Site, error) {
	var response struct {
		Data []*Site `json:"data"`
	}

	if err := u.GetData(APISiteList, &response); err != nil {
		return nil, err
	}

	sites := []string{} // used for debug log only

	for i, d := range response.Data {
		// Add the unifi struct to the site.
		response.Data[i].controller = u
		// Add special SourceName value.
		response.Data[i].SourceName = u.URL
		// If the human name is missing (description), set it to the cryptic name.
		response.Data[i].Desc = strings.TrimSpace(pick(d.Desc, d.Name))
		// Add the custom site name to each site. used as a Grafana filter somewhere.
		response.Data[i].SiteName = d.Desc + " (" + d.Name + ")"
		sites = append(sites, d.Name) // used for debug log only
	}

	u.DebugLog("Found %d site(s): %s", len(sites), strings.Join(sites, ","))

	return response.Data, nil
}

// GetSiteDPI garners dpi data for sites.
func (u *Unifi) GetSiteDPI(sites []*Site) ([]*DPITable, error) {
	data := []*DPITable{}

	for _, site := range sites {
		u.DebugLog("Polling Controller, retreiving Site DPI data, site %s", site.SiteName)

		var response struct {
			Data []*DPITable `json:"data"`
		}

		siteDPIpath := fmt.Sprintf(APISiteDPI, site.Name)
		if err := u.GetData(siteDPIpath, &response, `{"type":"by_app"}`); err != nil {
			return nil, err
		}

		if l := len(response.Data); l > 1 {
			return nil, ErrDPIDataBug
		} else if l == 0 {
			u.DebugLog("Site DPI data missing! Is DPI enabled in UniFi controller? Site %s", site.SiteName)
			continue
		}

		response.Data[0].SourceName = site.SourceName
		response.Data[0].SiteName = site.SiteName
		data = append(data, response.Data[0])
	}

	return data, nil
}

// Site represents a site's data.
type Site struct {
	controller   *Unifi
	SourceName   string   `json:"-"`
	ID           string   `json:"_id"`
	Name         string   `json:"name"`
	Desc         string   `json:"desc"`
	SiteName     string   `json:"-"`
	AttrHiddenID string   `json:"attr_hidden_id"`
	AttrNoDelete FlexBool `json:"attr_no_delete"`
	Health       []struct {
		Subsystem       string   `json:"subsystem"`
		NumUser         FlexInt  `json:"num_user,omitempty"`
		NumGuest        FlexInt  `json:"num_guest,omitempty"`
		NumIot          FlexInt  `json:"num_iot,omitempty"`
		TxBytesR        FlexInt  `json:"tx_bytes-r,omitempty"`
		RxBytesR        FlexInt  `json:"rx_bytes-r,omitempty"`
		Status          string   `json:"status"`
		NumAp           FlexInt  `json:"num_ap,omitempty"`
		NumAdopted      FlexInt  `json:"num_adopted,omitempty"`
		NumDisabled     FlexInt  `json:"num_disabled,omitempty"`
		NumDisconnected FlexInt  `json:"num_disconnected,omitempty"`
		NumPending      FlexInt  `json:"num_pending,omitempty"`
		NumGw           FlexInt  `json:"num_gw,omitempty"`
		WanIP           string   `json:"wan_ip,omitempty"`
		Gateways        []string `json:"gateways,omitempty"`
		Netmask         string   `json:"netmask,omitempty"`
		Nameservers     []string `json:"nameservers,omitempty"`
		NumSta          FlexInt  `json:"num_sta,omitempty"`
		GwMac           string   `json:"gw_mac,omitempty"`
		GwName          string   `json:"gw_name,omitempty"`
		GwSystemStats   struct {
			CPU    FlexInt `json:"cpu"`
			Mem    FlexInt `json:"mem"`
			Uptime FlexInt `json:"uptime"`
		} `json:"gw_system-stats,omitempty"`
		GwVersion             string   `json:"gw_version,omitempty"`
		Latency               FlexInt  `json:"latency,omitempty"`
		Uptime                FlexInt  `json:"uptime,omitempty"`
		Drops                 FlexInt  `json:"drops,omitempty"`
		XputUp                FlexInt  `json:"xput_up,omitempty"`
		XputDown              FlexInt  `json:"xput_down,omitempty"`
		SpeedtestStatus       string   `json:"speedtest_status,omitempty"`
		SpeedtestLastrun      FlexInt  `json:"speedtest_lastrun,omitempty"`
		SpeedtestPing         FlexInt  `json:"speedtest_ping,omitempty"`
		LanIP                 string   `json:"lan_ip,omitempty"`
		NumSw                 FlexInt  `json:"num_sw,omitempty"`
		RemoteUserEnabled     FlexBool `json:"remote_user_enabled,omitempty"`
		RemoteUserNumActive   FlexInt  `json:"remote_user_num_active,omitempty"`
		RemoteUserNumInactive FlexInt  `json:"remote_user_num_inactive,omitempty"`
		RemoteUserRxBytes     FlexInt  `json:"remote_user_rx_bytes,omitempty"`
		RemoteUserTxBytes     FlexInt  `json:"remote_user_tx_bytes,omitempty"`
		RemoteUserRxPackets   FlexInt  `json:"remote_user_rx_packets,omitempty"`
		RemoteUserTxPackets   FlexInt  `json:"remote_user_tx_packets,omitempty"`
		SiteToSiteEnabled     FlexBool `json:"site_to_site_enabled,omitempty"`
	} `json:"health"`
	NumNewAlarms FlexInt `json:"num_new_alarms"`
}
