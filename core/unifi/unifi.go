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
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

var (
	errAuthenticationFailed = fmt.Errorf("authentication failed")
	errInvalidStatusCode    = fmt.Errorf("invalid status code from server")
)

// NewUnifi creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
// Start here.
func NewUnifi(config *Config) (*Unifi, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
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
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.VerifySSL}, // nolint: gosec
			},
		},
	}

	if err := u.checkNewStyleAPI(); err != nil {
		return u, err
	}

	if err := u.Login(); err != nil {
		return u, err
	}

	if err := u.GetServerData(); err != nil {
		return u, errors.Wrap(err, "unable to get server version")
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

	defer resp.Body.Close()                   // we need no data here.
	_, _ = io.Copy(ioutil.Discard, resp.Body) // avoid leaking.
	u.DebugLog("Requested %s: elapsed %v, returned %d bytes",
		req.URL, time.Since(start).Round(time.Millisecond), resp.ContentLength)

	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(errAuthenticationFailed, "(user: %s): %s (status: %s)",
			u.User, req.URL, resp.Status)
	}

	return nil
}

// with the release of controller version 5.12.55 on UDM in Jan 2020 the api paths
// changed and broke this library. This function runs when `NewUnifi()` is called to
// check if this is a newer controller or not. If it is, we set new to true.
// Setting new to true makes the path() method return different (new) paths.
func (u *Unifi) checkNewStyleAPI() error {
	u.DebugLog("Requesting %s/ to determine API paths", u.URL)

	req, err := http.NewRequest("GET", u.URL+"/", nil)
	if err != nil {
		return err
	}

	// We can't share these cookies with other requests, so make a new client.
	// Checking the return code on the first request so don't follow a redirect.
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !u.VerifySSL}, // nolint: gosec
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()                   // we need no data here.
	_, _ = io.Copy(ioutil.Discard, resp.Body) // avoid leaking.

	if resp.StatusCode == http.StatusOK {
		// The new version returns a "200" for a / request.
		u.New = true
		u.DebugLog("Using NEW UniFi controller API paths for %s", req.URL)
	}

	// The old version returns a "302" (to /manage) for a / request
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
func (u *Unifi) GetData(apiPath string, v interface{}, params ...string) error {
	start := time.Now()

	body, err := u.GetJSON(apiPath, params...)
	if err != nil {
		return err
	}

	u.DebugLog("Requested %s: elapsed %v, returned %d bytes",
		u.URL+u.path(apiPath), time.Since(start).Round(time.Millisecond), len(body))

	return json.Unmarshal(body, v)
}

// UniReq is a small helper function that adds an Accept header.
// Use this if you're unmarshalling UniFi data into custom types.
// And if you're doing that... sumbut a pull request with your new struct. :)
// This is a helper method that is exposed for convenience.
func (u *Unifi) UniReq(apiPath string, params string) (req *http.Request, err error) {
	switch apiPath = u.path(apiPath); params {
	case "":
		req, err = http.NewRequest("GET", u.URL+apiPath, nil)
	default:
		req, err = http.NewRequest("POST", u.URL+apiPath, bytes.NewBufferString(params))
	}

	if err != nil {
		return
	}

	// Add the saved CSRF header.
	req.Header.Set("X-CSRF-Token", u.csrf)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if u.Client.Jar != nil {
		parsedURL, _ := url.Parse(req.URL.String())
		u.DebugLog("Requesting %s, with params: %v, cookies: %d", req.URL, params != "", len(u.Client.Jar.Cookies(parsedURL)))
	} else {
		u.DebugLog("Requesting %s, with params: %v,", req.URL, params != "")
	}

	return
}

// GetJSON returns the raw JSON from a path. This is useful for debugging.
func (u *Unifi) GetJSON(apiPath string, params ...string) ([]byte, error) {
	req, err := u.UniReq(apiPath, strings.Join(params, " "))
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

	// Save the returned CSRF header.
	if csrf := resp.Header.Get("x-csrf-token"); csrf != "" {
		u.csrf = resp.Header.Get("x-csrf-token")
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Wrapf(errInvalidStatusCode, "%s: %s", req.URL, resp.Status)
	}

	return body, err
}
