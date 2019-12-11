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
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

// NewUnifi creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
// Start here.
func NewUnifi(config *Config) (*Unifi, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	config.URL = strings.TrimRight(config.URL, "/")

	if config.ErrorLog == nil {
		config.ErrorLog = discardLogs
	}

	if config.DebugLog == nil {
		config.DebugLog = discardLogs
	}

	u := &Unifi{Config: config,
		Client: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.VerifySSL},
			},
		},
	}
	if err := u.Login(); err != nil {
		return u, err
	}

	if err := u.GetServerData(); err != nil {
		return u, fmt.Errorf("unable to get server version: %v", err)
	}

	return u, nil
}

// Login is a helper method. It can be called to grab a new authentication cookie.
func (u *Unifi) Login() error {
	start := time.Now()

	// magic login.
	req, err := u.UniReq(APILoginPath, fmt.Sprintf(`{"username":"%s","password":"%s"}`, u.User, u.Pass))
	if err != nil {
		return err
	}

	resp, err := u.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // we need no data here.

	_, _ = io.Copy(ioutil.Discard, resp.Body) // avoid leaking.
	u.DebugLog("Requested %s: elapsed %v, returned %d bytes",
		APILoginPath, time.Since(start).Round(time.Millisecond), resp.ContentLength)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed (user: %s): %s (status: %s)",
			u.User, u.URL+APILoginPath, resp.Status)
	}

	return nil
}

// GetServerData sets the controller's version and UUID. Only call this if you
// previously called Login and suspect the controller version has changed.
func (u *Unifi) GetServerData() error {
	var response struct {
		Data server `json:"meta"`
	}

	u.server = &response.Data

	return u.GetData(APIStatusPath, &response)
}

// GetData makes a unifi request and unmarshals the response into a provided pointer.
func (u *Unifi) GetData(apiPath string, v interface{}) error {
	start := time.Now()

	body, err := u.GetJSON(apiPath)
	if err != nil {
		return err
	}

	u.DebugLog("Requested %s: elapsed %v, returned %d bytes",
		apiPath, time.Since(start).Round(time.Millisecond), len(body))

	return json.Unmarshal(body, v)
}

// UniReq is a small helper function that adds an Accept header.
// Use this if you're unmarshalling UniFi data into custom types.
// And if you're doing that... sumbut a pull request with your new struct. :)
// This is a helper method that is exposed for convenience.
func (u *Unifi) UniReq(apiPath string, params string) (req *http.Request, err error) {
	switch params {
	case "":
		req, err = http.NewRequest("GET", u.URL+apiPath, nil)
	default:
		req, err = http.NewRequest("POST", u.URL+apiPath, bytes.NewBufferString(params))
	}

	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	u.DebugLog("Requesting %s, with params: %v", apiPath, params != "")

	return
}

// GetJSON returns the raw JSON from a path. This is useful for debugging.
func (u *Unifi) GetJSON(apiPath string) ([]byte, error) {
	req, err := u.UniReq(apiPath, "")
	if err != nil {
		return []byte{}, err
	}

	resp, err := u.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("invalid status code from server %s", resp.Status)
	}

	return body, err
}
