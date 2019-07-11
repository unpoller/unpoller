// Package unifi provides a set of types to unload (unmarshal) Ubiquiti UniFi
// controller data. Also provided are methods to easily get data for devices -
// things like access points and switches, and for clients - the things
// connected to those access points and switches. As a bonus, each device and
// client type provided has an attached method to create InfluxDB datapoints.
package unifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/pkg/errors"
)

// NewUnifi creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
// Start here.
func NewUnifi(user, pass, url string, verifySSL bool) (*Unifi, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "cookiejar.New(nil)")
	}
	u := &Unifi{baseURL: strings.TrimRight(url, "/"),
		Client: &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL}},
			Jar:       jar,
		},
		ErrorLog: log.Printf,
		DebugLog: DiscardLogs,
	}
	return u, u.getController(user, pass)
}

// getController is a helper method to make testsing a bit easier.
func (u *Unifi) getController(user, pass string) error {
	// magic login.
	req, err := u.UniReq(LoginPath, fmt.Sprintf(`{"username":"%s","password":"%s"}`, user, pass))
	if err != nil {
		return errors.Wrap(err, "UniReq(LoginPath, json)")
	}
	resp, err := u.Do(req)
	if err != nil {
		return errors.Wrap(err, "authReq.Do(req)")
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body) // avoid leaking.
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("authentication failed (user: %s): %s (status: %s)",
			user, u.baseURL+LoginPath, resp.Status)
	}
	return errors.Wrap(u.getServer(), "unable to get server version")
}

// getServer sets the controller's version and UUID.
func (u *Unifi) getServer() error {
	var response struct {
		Data server `json:"meta"`
	}
	u.server = &response.Data
	return u.GetData(StatusPath, &response)
}

// GetClients returns a response full of clients' data from the UniFi Controller.
func (u *Unifi) GetClients(sites []Site) (Clients, error) {
	data := make([]Client, 0)
	for _, site := range sites {
		var response struct {
			Data []Client `json:"data"`
		}
		u.DebugLog("Polling Controller, retreiving UniFi Clients, site %s (%s) ", site.Name, site.Desc)
		clientPath := fmt.Sprintf(ClientPath, site.Name)
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}
		for i, d := range response.Data {
			// Add the special "Site Name" to each client. This becomes a Grafana filter somewhere.
			response.Data[i].SiteName = site.Desc + " (" + site.Name + ")"
			// Fix name and hostname fields. Sometimes one or the other is blank.
			response.Data[i].Hostname = pick(d.Hostname, d.Name, d.Mac)
			response.Data[i].Name = pick(d.Name, d.Hostname)
		}
		data = append(data, response.Data...)
	}
	return data, nil
}

// GetDevices returns a response full of devices' data from the UniFi Controller.
func (u *Unifi) GetDevices(sites []Site) (*Devices, error) {
	devices := new(Devices)
	for _, site := range sites {
		var response struct {
			Data []json.RawMessage `json:"data"`
		}
		devicePath := fmt.Sprintf(DevicePath, site.Name)
		if err := u.GetData(devicePath, &response); err != nil {
			return nil, err
		}
		loopDevices := u.parseDevices(response.Data, site.SiteName)
		devices.UAPs = append(devices.UAPs, loopDevices.UAPs...)
		devices.USGs = append(devices.USGs, loopDevices.USGs...)
		devices.USWs = append(devices.USWs, loopDevices.USWs...)
	}
	return devices, nil
}

// GetSites returns a list of configured sites on the UniFi controller.
func (u *Unifi) GetSites() (Sites, error) {
	var response struct {
		Data []Site `json:"data"`
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

// GetData makes a unifi request and unmarshal the response into a provided pointer.
func (u *Unifi) GetData(methodPath string, v interface{}) error {
	body, err := u.GetJSON(methodPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	return errors.Wrapf(err, "json.Unmarshal(%s)", methodPath)
}

// UniReq is a small helper function that adds an Accept header.
// Use this if you're unmarshalling UniFi data into custom types.
// And if you're doing that... sumbut a pull request with your new struct. :)
// This is a helper method that is exposed for convenience.
func (u *Unifi) UniReq(apiPath string, params string) (req *http.Request, err error) {
	switch path := u.baseURL + apiPath; {
	case params == "":
		req, err = http.NewRequest("GET", path, nil)
	default:
		req, err = http.NewRequest("POST", path, bytes.NewBufferString(params))
	}
	if err == nil {
		req.Header.Add("Accept", "application/json")
	}
	return
}

// GetJSON returns the raw JSON from a path. This is useful for debugging.
func (u *Unifi) GetJSON(apiPath string) ([]byte, error) {
	req, err := u.UniReq(apiPath, "")
	if err != nil {
		return []byte{}, errors.Wrapf(err, "c.UniReq(%s)", apiPath)
	}
	resp, err := u.Do(req)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "c.Do(%s)", apiPath)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, errors.Wrapf(err, "ioutil.ReadAll(%s)", apiPath)
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("invalid status code from server %s", resp.Status)
	}
	return body, err
}
