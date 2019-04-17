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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/pkg/errors"
)

// GetController creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
// Start here.
func GetController(user, pass, url string, verifySSL bool) (*Unifi, error) {
	u := &Unifi{baseURL: strings.TrimRight(url, "/")}
	json := `{"username": "` + user + `","password": "` + pass + `"}`
	req, err := u.UniReq(LoginPath, json)
	if err != nil {
		return u, errors.Wrap(err, "UniReq(LoginPath, json)")
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "cookiejar.New(nil)")
	}
	u.Client = &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL}},
		Jar:       jar,
	}
	return u, u.getController(req)
}

// getController is a helper method to make testsing a bit easier.
func (u *Unifi) getController(req *http.Request) error {
	resp, err := u.Do(req)
	if err != nil {
		return errors.Wrap(err, "authReq.Do(req)")
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body) // avoid leaking.
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("authentication failed: %v (status: %v/%v)",
			u.baseURL+LoginPath, resp.StatusCode, resp.Status)
	}
	return nil
}

// GetClients returns a response full of clients' data from the Unifi Controller.
func (u *Unifi) GetClients() (*Clients, error) {
	var response struct {
		Clients []UCL `json:"data"`
	}
	req, err := u.UniReq(ClientPath, "")
	if err != nil {
		return nil, errors.Wrap(err, "c.UniReq(ClientPath)")
	}
	resp, err := u.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "c.Do(req)")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll(resp.Body)")
	} else if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal([]UCL)")
	}
	return &Clients{UCLs: response.Clients}, nil
}

// GetDevices returns a response full of devices' data from the Unifi Controller.
func (u *Unifi) GetDevices() (*Devices, error) {
	var response struct {
		Data []json.RawMessage `json:"data"`
	}
	req, err := u.UniReq(DevicePath, "")
	if err != nil {
		return nil, errors.Wrap(err, "c.UniReq(DevicePath)")
	}
	resp, err := u.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "c.Do(req)")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadAll(resp.Body)")
	} else if err = json.Unmarshal(body, &response); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal([]json.RawMessage)")
	}
	return u.parseDevices(response.Data), nil
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
