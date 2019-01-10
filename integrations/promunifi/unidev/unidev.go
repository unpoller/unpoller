package unidev

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
	"net/http/cookiejar"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
)

// LoginPath is Unifi Controller Login API Path
const LoginPath = "/api/login"

// Asset provides a common interface to retreive metrics from a device or client.
// It currently only supports InfluxDB, but could be amended to support other
// libraries that have a similar interface.
// This app only uses the .AddPoint/s() methods with the Asset type.
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

// StringInt is used to unmarshal quoted integers in JSON responses.
type StringInt int

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
	req, err := authReq.UniReq(LoginPath, json)
	if err != nil {
		return nil, errors.Wrap(err, "UniReq(LoginPath, json)")
	}
	resp, err := authReq.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "authReq.Do(req)")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("resp.Body.Close():", err) // Not fatal. Just log it.
		}
	}()
	if resp.StatusCode != http.StatusOK {
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
