package unidev

import (
	"bytes"
	"crypto/tls"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"net/http"
	"net/http/cookiejar"
)

// LoginPath is Unifi Controller Login API Path
const LoginPath = "/api/login"

// Asset provides a common interface to retreive metrics from a device or client.
type Asset interface {
	// Point() means this is useful to influxdb..
	Points() ([]*influx.Point, error)
	// Add more methods to achieve more usefulness from this library.
}

// AuthedReq is what you get in return for providing a password!
type AuthedReq struct {
	*http.Client
	baseURL string
}

// AuthController creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
func AuthController(user, pass, url string) (*AuthedReq, error) {
	json := `{"username": "` + user + `","password": "` + pass + `"}`
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "cookiejar.New(nil)")
	}
	authReq := &AuthedReq{&http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Jar:       jar,
	}, url}
	if req, err := authReq.UniReq(LoginPath, json); err != nil {
		return nil, errors.Wrap(err, "uniRequest(LoginPath, json)")
	} else if resp, err := authReq.Do(req); err != nil {
		return nil, errors.Wrap(err, "c.uniClient.Do(req)")
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("authentication failed (%v): %v (status: %v/%v)",
			user, url+LoginPath, resp.StatusCode, resp.Status)
	}
	return authReq, nil
}

// UniReq is a small helper function that adds an Accept header.
func (c AuthedReq) UniReq(apiURL string, params string) (req *http.Request, err error) {
	if params != "" {
		req, err = http.NewRequest("POST", c.baseURL+apiURL, bytes.NewBufferString(params))
	} else {
		req, err = http.NewRequest("GET", c.baseURL+apiURL, nil)
	}
	if err == nil {
		req.Header.Add("Accept", "application/json")
	}
	return
}
