package unifi

import "strings"

// GetSites returns a list of configured sites on the UniFi controller.
func (u *Unifi) GetSites() (Sites, error) {
	var response struct {
		Data []*Site `json:"data"`
	}
	if err := u.GetData(SiteList, &response); err != nil {
		return nil, err
	}
	sites := []string{} // used for debug log only
	for i, d := range response.Data {
		// If the human name is missing (description), set it to the cryptic name.
		response.Data[i].Desc = pick(d.Desc, d.Name)
		// Add the custom site name to each site. used as a Grafana filter somewhere.
		response.Data[i].SiteName = d.Desc + " (" + d.Name + ")"
		sites = append(sites, d.Name) // used for debug log only
	}
	u.DebugLog("Found %d site(s): %s", len(sites), strings.Join(sites, ","))
	return response.Data, nil
}

// Sites is a struct to match Devices and Clients.
type Sites []*Site

// Site represents a site's data.
type Site struct {
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
