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

// Logger is a base type to deal with changing log outs.
type Logger func(msg string, fmt ...interface{})

// LoginPath is Unifi Controller Login API Path
const LoginPath = "/api/login"

// Devices contains a list of all the unifi devices from a controller.
type Devices struct {
	UAPs []UAP
	USGs []USG
	USWs []USW
}

// AuthedReq is what you get in return for providing a password!
type AuthedReq struct {
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

// AuthController creates a http.Client with authenticated cookies.
// Used to make additional, authenticated requests to the APIs.
func AuthController(user, pass, url string, verifySSL bool) (*AuthedReq, error) {
	json := `{"username": "` + user + `","password": "` + pass + `"}`
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "cookiejar.New(nil)")
	}
	a := &AuthedReq{Client: &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL}},
		Jar:       jar,
	}, baseURL: url}
	req, err := a.UniReq(LoginPath, json)
	if err != nil {
		return a, errors.Wrap(err, "UniReq(LoginPath, json)")
	}
	resp, err := a.Do(req)
	if err != nil {
		return a, errors.Wrap(err, "authReq.Do(req)")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return a, errors.Errorf("authentication failed (%v): %v (status: %v/%v)",
			user, url+LoginPath, resp.StatusCode, resp.Status)
	}
	return a, nil
}

// UniReq is a small helper function that adds an Accept header.
func (a AuthedReq) UniReq(apiPath string, params string) (req *http.Request, err error) {
	if params != "" {
		req, err = http.NewRequest("POST", a.baseURL+apiPath, bytes.NewBufferString(params))
	} else {
		req, err = http.NewRequest("GET", a.baseURL+apiPath, nil)
	}
	if err == nil {
		req.Header.Add("Accept", "application/json")
	}
	return
}

func (a AuthedReq) dLogf(msg string, v ...interface{}) {
	if a.DebugLog != nil {
		a.DebugLog("[DEBUG] "+msg, v...)
	}
}

func (a AuthedReq) eLogf(msg string, v ...interface{}) {
	if a.ErrorLog != nil {
		a.ErrorLog("[ERROR] "+msg, v...)
	}
}
