// Package unifi provides a set of types to unload (unmarshal) Unifi Ubiquiti
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
	}
	return u, u.getController(user, pass)
}

// getController is a helper method to make testsing a bit easier.
func (u *Unifi) getController(user, pass string) error {
	// magic login.
	req, err := u.UniReq(LoginPath, `{"username": "`+user+`","password": "`+pass+`"}`)
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
		return errors.Errorf("authentication failed (user: %s): %s (status: %d/%s)",
			user, u.baseURL+LoginPath, resp.StatusCode, resp.Status)
	}
	return nil
}

// GetServer returns the controller's version and UUID.
func (u *Unifi) GetServer() (Server, error) {
	var response struct {
		Data Server `json:"meta"`
	}
	err := u.GetData(StatusPath, &response)
	return response.Data, err
}

// GetClients returns a response full of clients' data from the Unifi Controller.
func (u *Unifi) GetClients(sites []Site) (Clients, error) {
	data := make([]Client, 0)
	for _, site := range sites {
		var response struct {
			Data []Client `json:"data"`
		}
		u.dLogf("Polling Controller, retreiving Unifi Clients, site %s (%s) ", site.Name, site.Desc)
		clientPath := fmt.Sprintf(ClientPath, site.Name)
		if err := u.GetData(clientPath, &response); err != nil {
			return nil, err
		}
		for i := range response.Data {
			response.Data[i].SiteName = site.Desc + " (" + site.Name + ")"
		}
		data = append(data, response.Data...)
	}
	return data, nil
}

// GetDevices returns a response full of devices' data from the Unifi Controller.
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
		loopDevices := u.parseDevices(response.Data, site.Desc+" ("+site.Name+")")
		// Add SiteName to each device asset.
		for i := range loopDevices.UAPs {
			loopDevices.UAPs[i].SiteName = site.Desc + " (" + site.Name + ")"
			for j := range loopDevices.UAPs[i].VapTable {
				loopDevices.UAPs[i].VapTable[j].SiteName = site.Desc + " (" + site.Name + ")"
			}
		}
		for i := range loopDevices.USGs {
			loopDevices.USGs[i].SiteName = site.Desc + " (" + site.Name + ")"
			for j := range loopDevices.USGs[i].NetworkTable {
				loopDevices.USGs[i].NetworkTable[j].SiteName = site.Desc + " (" + site.Name + ")"
			}
		}
		for i := range loopDevices.USWs {
			loopDevices.USWs[i].SiteName = site.Desc + " (" + site.Name + ")"
		}
		devices.UAPs = append(devices.UAPs, loopDevices.UAPs...)
		devices.USGs = append(devices.USGs, loopDevices.USGs...)
		devices.USWs = append(devices.USWs, loopDevices.USWs...)
	}
	return devices, nil
}

// GetSites returns a list of configured sites on the Unifi controller.
func (u *Unifi) GetSites() (Sites, error) {
	var response struct {
		Data []Site `json:"data"`
	}
	if err := u.GetData(SiteList, &response); err != nil {
		return nil, err
	}
	sites := make([]string, 0)
	for i := range response.Data {
		response.Data[i].SiteName = response.Data[i].Desc + " (" + response.Data[i].Name + ")"
		sites = append(sites, response.Data[i].Name)
	}
	u.dLogf("Found %d site(s): %s", len(sites), strings.Join(sites, ","))
	return response.Data, nil
}

// GetData makes a unifi request and unmarshal the response into a provided pointer.
func (u *Unifi) GetData(methodPath string, v interface{}) error {
	req, err := u.UniReq(methodPath, "")
	if err != nil {
		return errors.Wrapf(err, "c.UniReq(%s)", methodPath)
	}
	resp, err := u.Do(req)
	if err != nil {
		return errors.Wrapf(err, "c.Do(%s)", methodPath)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return errors.Wrapf(err, "ioutil.ReadAll(%s)", methodPath)
	} else if err = json.Unmarshal(body, v); err != nil {
		return errors.Wrapf(err, "json.Unmarshal(%s)", methodPath)
	}
	return nil
}

// UniReq is a small helper function that adds an Accept header.
// Use this if you're unmarshalling Unifi data into custom types.
// And if you're doing that... sumbut a pull request with your new struct. :)
// This is a helper method that is exposed for convenience.
func (u *Unifi) UniReq(apiPath string, params string) (req *http.Request, err error) {
	if params != "" {
		req, err = http.NewRequest("POST", u.baseURL+apiPath, bytes.NewBufferString(params))
	} else {
		req, err = http.NewRequest("GET", u.baseURL+apiPath, nil)
	}
	if err == nil {
		req.Header.Add("Accept", "application/json")
	}
	return
}
