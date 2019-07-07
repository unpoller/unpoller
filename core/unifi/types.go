package unifi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// This is a list of unifi API paths.
// The %s in each string must be replaced with a Site.Name.
const (
	// StatusPath shows Controller version.
	StatusPath string = "/status"
	// SiteList is the path to the api site list.
	SiteList string = "/api/stat/sites"
	// ClientPath is Unifi Clients API Path
	ClientPath string = "/api/s/%s/stat/sta"
	// DevicePath is where we get data about Unifi devices.
	DevicePath string = "/api/s/%s/stat/device"
	// NetworkPath contains network-configuration data. Not really graphable.
	NetworkPath string = "/api/s/%s/rest/networkconf"
	// UserGroupPath contains usergroup configurations.
	UserGroupPath string = "/api/s/%s/rest/usergroup"
	// LoginPath is Unifi Controller Login API Path
	LoginPath string = "/api/login"
)

// Logger is a base type to deal with changing log outputs. Create a logger
// that matches this interface to capture debug and error logs.
type Logger func(msg string, fmt ...interface{})

// DiscardLogs is the default debug logger.
func DiscardLogs(msg string, v ...interface{}) {
	// do nothing.
}

// Devices contains a list of all the unifi devices from a controller.
// Contains Access points, security gateways and switches.
type Devices struct {
	UAPs []UAP
	USGs []USG
	USWs []USW
}

// Unifi is what you get in return for providing a password! Unifi represents
// a controller that you can make authenticated requests to. Use this to make
// additional requests for devices, clients or other custom data. Do not set
// the loggers to nil. Set them to DiscardLogs if you want no logs.
type Unifi struct {
	*http.Client
	baseURL  string
	ErrorLog Logger
	DebugLog Logger
	*server
}

// server is the /status endpoint from the Unifi controller.
type server struct {
	Up            FlexBool `json:"up"`
	ServerVersion string   `json:"server_version"`
	UUID          string   `json:"uuid"`
}

// FlexInt provides a container and unmarshalling for fields that may be
// numbers or strings in the Unifi API.
type FlexInt struct {
	Val float64
	Txt string
}

// UnmarshalJSON converts a string or number to an integer.
// Generally, do call this directly, it's used in the json interface.
func (f *FlexInt) UnmarshalJSON(b []byte) error {
	var unk interface{}
	if err := json.Unmarshal(b, &unk); err != nil {
		return err
	}
	switch i := unk.(type) {
	case float64:
		f.Val = i
		f.Txt = strconv.FormatFloat(i, 'f', -1, 64)
		return nil
	case string:
		f.Txt = i
		f.Val, _ = strconv.ParseFloat(i, 64)
		return nil
	default:
		return errors.New("Cannot unmarshal to FlexInt")
	}
}

// FlexBool provides a container and unmarshalling for fields that may be
// boolean or strings in the Unifi API.
type FlexBool struct {
	Val bool
	Txt string
}

// UnmarshalJSON method converts armed/disarmed, yes/no, active/inactive or 0/1 to true/false.
// Really it converts ready, ok, up, t, armed, yes, active, enabled, 1, true to true. Anything else is false.
func (f *FlexBool) UnmarshalJSON(b []byte) error {
	if f.Txt = strings.Trim(string(b), `"`); f.Txt == "" {
		f.Txt = "false"
	}
	f.Val = f.Txt == "1" || strings.EqualFold(f.Txt, "true") || strings.EqualFold(f.Txt, "yes") ||
		strings.EqualFold(f.Txt, "t") || strings.EqualFold(f.Txt, "armed") || strings.EqualFold(f.Txt, "active") ||
		strings.EqualFold(f.Txt, "enabled") || strings.EqualFold(f.Txt, "ready") || strings.EqualFold(f.Txt, "up") ||
		strings.EqualFold(f.Txt, "ok")
	return nil
}
