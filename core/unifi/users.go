package unifi

import (
	"fmt"
	"strings"
)

// GetUsers returns a response full of clients that connected to the UDM within the provided amount of time
// using the insight historical connection data set.
func (u *Unifi) GetUsers(sites []*Site, hours int) ([]*User, error) {
	data := make([]*User, 0)

	for _, site := range sites {
		var (
			response struct {
				Data []*User `json:"data"`
			}
			params = fmt.Sprintf(`{ "type": "all:", "conn": "all", "within":%d }`, hours)
		)

		u.DebugLog("Polling Controller, retrieving UniFi Users, site %s ", site.SiteName)

		clientPath := fmt.Sprintf(APIAllUserPath, site.Name)
		if err := u.GetData(clientPath, &response, params); err != nil {
			return nil, err
		}

		for i, d := range response.Data {
			// Add special SourceName value.
			response.Data[i].SourceName = u.URL
			// Add the special "Site Name" to each client. This becomes a Grafana filter somewhere.
			response.Data[i].SiteName = site.SiteName
			// Fix name and hostname fields. Sometimes one or the other is blank.
			response.Data[i].Hostname = strings.TrimSpace(pick(d.Hostname, d.Name, d.Mac))
			response.Data[i].Name = strings.TrimSpace(pick(d.Name, d.Hostname))
		}

		data = append(data, response.Data...)
	}

	return data, nil
}

// User defines the metadata available for previously connected clients.
type User struct {
	SourceName          string   `json:"-"`
	SiteName            string   `json:"-"`
	ID                  string   `json:"_id"`
	Mac                 string   `json:"mac"`
	SiteID              string   `json:"site_id"`
	Oui                 string   `json:"oui,omitempty"`
	IsGuest             bool     `json:"is_guest"`
	FirstSeen           FlexInt  `json:"first_seen,omitempty"`
	LastSeen            FlexInt  `json:"last_seen,omitempty"`
	IsWired             bool     `json:"is_wired,omitempty"`
	Hostname            string   `json:"hostname,omitempty"`
	Duration            FlexInt  `json:"duration,omitempty"`
	TxBytes             FlexInt  `json:"tx_bytes,omitempty"`
	TxPackets           FlexInt  `json:"tx_packets,omitempty"`
	RxBytes             FlexInt  `json:"rx_bytes,omitempty"`
	RxPackets           FlexInt  `json:"rx_packets,omitempty"`
	WifiTxAttempts      FlexInt  `json:"wifi_tx_attempts,omitempty"`
	TxRetries           FlexInt  `json:"tx_retries,omitempty"`
	UsergroupID         string   `json:"usergroup_id,omitempty"`
	Name                string   `json:"name,omitempty"`
	Note                string   `json:"note,omitempty"`
	Noted               FlexBool `json:"noted,omitempty"`
	Blocked             FlexBool `json:"blocked,omitempty"`
	DevIDOverride       FlexInt  `json:"dev_id_override,omitempty"`
	FingerprintOverride FlexBool `json:"fingerprint_override,omitempty"`
}
