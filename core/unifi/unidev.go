package unifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"strconv"

	"github.com/pkg/errors"
)

const (
	// ClientPath is Unifi Clients API Path
	ClientPath string = "/api/s/default/stat/sta"
	// DevicePath is where we get data about Unifi devices.
	DevicePath string = "/api/s/default/stat/device"
	// NetworkPath contains network-configuration data. Not really graphable.
	NetworkPath string = "/api/s/default/rest/networkconf"
	// UserGroupPath contains usergroup configurations.
	UserGroupPath string = "/api/s/default/rest/usergroup"
	// LoginPath is Unifi Controller Login API Path
	LoginPath string = "/api/login"
)

// Logger is a base type to deal with changing log outputs.
type Logger func(msg string, fmt ...interface{})

// Devices contains a list of all the unifi devices from a controller.
type Devices struct {
	UAPs []UAP
	USGs []USG
	USWs []USW
}

// Clients conptains a list of all the unifi clients from a controller.
type Clients struct {
	UCLs []UCL
}

// Unifi is what you get in return for providing a password!
type Unifi struct {
	*http.Client
	baseURL  string
	ErrorLog Logger
	DebugLog Logger
}

// FlexInt provides a container and unmarshalling for fields that may be
// numbers or strings in the Unifi API
type FlexInt struct {
	Number float64
	String string
}

// UnmarshalJSON converts a string or number to an integer.
func (f *FlexInt) UnmarshalJSON(b []byte) error {
	var unk interface{}
	if err := json.Unmarshal(b, &unk); err != nil {
		return err
	}
	switch i := unk.(type) {
	case float64:
		f.Number = i
		f.String = strconv.FormatFloat(i, 'f', -1, 64)
		return nil
	case string:
		f.String = i
		f.Number, _ = strconv.ParseFloat(i, 64)
		return nil
	default:
		return errors.New("Cannot unmarshal to FlexInt")
	}
}

// GetController creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
func GetController(user, pass, url string, verifySSL bool) (*Unifi, error) {
	json := `{"username": "` + user + `","password": "` + pass + `"}`
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "cookiejar.New(nil)")
	}
	u := &Unifi{
		Client: &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL}},
			Jar:       jar,
		},
		baseURL: url,
	}
	req, err := u.UniReq(LoginPath, json)
	if err != nil {
		return u, errors.Wrap(err, "UniReq(LoginPath, json)")
	}
	resp, err := u.Do(req)
	if err != nil {
		return u, errors.Wrap(err, "authReq.Do(req)")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return u, errors.Errorf("authentication failed (%v): %v (status: %v/%v)",
			user, url+LoginPath, resp.StatusCode, resp.Status)
	}
	return u, nil
}

// UniReq is a small helper function that adds an Accept header.
// Use this if you're unmarshalling Unifi data into custom types.
// And you're doing that... sumbut a pull request with your new struct. :)
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

// dLogf logs a debug message.
func (u *Unifi) dLogf(msg string, v ...interface{}) {
	if u.DebugLog != nil {
		u.DebugLog("[DEBUG] "+msg, v...)
	}
}

// dLogf logs an error message.
func (u *Unifi) eLogf(msg string, v ...interface{}) {
	if u.ErrorLog != nil {
		u.ErrorLog("[ERROR] "+msg, v...)
	}
}
