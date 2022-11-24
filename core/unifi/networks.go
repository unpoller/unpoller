package unifi

import (
	"encoding/json"
	"fmt"
)

// GetNetworks returns a response full of network data from the UniFi Controller.
func (u *Unifi) GetNetworks(sites []*Site) ([]Network, error) {
	networks := make([]Network, 0)

	for _, site := range sites {
		var response struct {
			Data []json.RawMessage `json:"data"`
		}

		networkPath := fmt.Sprintf(APINetworkPath, site.Name)
		if err := u.GetData(networkPath, &response); err != nil {
			return nil, err
		}

		for _, data := range response.Data {
			network, err := u.parseNetwork(data, site.SiteName)
			if err != nil {
				return networks, err
			}

			networks = append(networks, *network)
		}
	}

	return networks, nil
}

// parseNetwork parses the raw JSON from the Unifi Controller into network structures.
func (u *Unifi) parseNetwork(data json.RawMessage, siteName string) (*Network, error) {
	network := new(Network)
	return network, u.unmarshalDevice(siteName, data, network)
}

// Network is metadata about a network managed by a UniFi controller.
type Network struct {
	DhcpdDNSEnabled        FlexBool `json:"dhcpd_dns_enabled"`
	DhcpdEnabled           FlexBool `json:"dhcpd_enabled"`
	DhcpdGatewayEnabled    FlexBool `json:"dhcpd_gateway_enabled"`
	DhcpdIP1               string   `json:"dhcpd_ip_1"`
	DhcpdLeasetime         FlexInt  `json:"dhcpd_leasetime"`
	DhcpRelayEnabled       FlexBool `json:"dhcp_relay_enabled"`
	DhcpdTimeOffsetEnabled FlexBool `json:"dhcpd_time_offset_enabled"`
	DhcpGuardEnabled       FlexBool `json:"dhcpguard_enabled"`
	DomainName             string   `json:"domain_name"`
	Enabled                FlexBool `json:"enabled"`
	ID                     string   `json:"_id"`
	IPSubnet               string   `json:"ip_subnet"`
	IsNat                  FlexBool `json:"is_nat"`
	Name                   string   `json:"name"`
	Networkgroup           string   `json:"networkgroup"`
	Purpose                string   `json:"purpose"`
	SiteID                 string   `json:"site_id"`
	Vlan                   FlexInt  `json:"vlan"`
	VlanEnabled            FlexBool `json:"vlan_enabled"`
}
