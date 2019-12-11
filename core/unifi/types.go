package unifi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// This is a list of unifi API paths.
// The %s in each string must be replaced with a Site.Name.
const (
	// APIStatusPath shows Controller version.
	APIStatusPath string = "/status"
	// APISiteList is the path to the api site list.
	APISiteList string = "/api/stat/sites"
	// APIClientPath is Unifi Clients API Path
	APIClientPath string = "/api/s/%s/stat/sta"
	// APIDevicePath is where we get data about Unifi devices.
	APIDevicePath string = "/api/s/%s/stat/device"
	// APINetworkPath contains network-configuration data. Not really graphable.
	APINetworkPath string = "/api/s/%s/rest/networkconf"
	// APIUserGroupPath contains usergroup configurations.
	APIUserGroupPath string = "/api/s/%s/rest/usergroup"
	// APILoginPath is Unifi Controller Login API Path
	APILoginPath string = "/api/login"
	// APIIPSEvents returns Intrusion Detection Systems Events
	APIIPSEvents string = "/api/s/%s/stat/ips/event"
)

// Logger is a base type to deal with changing log outputs. Create a logger
// that matches this interface to capture debug and error logs.
type Logger func(msg string, fmt ...interface{})

// discardLogs is the default debug logger.
func discardLogs(msg string, v ...interface{}) {
	// do nothing.
}

// Devices contains a list of all the unifi devices from a controller.
// Contains Access points, security gateways and switches.
type Devices struct {
	UAPs []*UAP
	USGs []*USG
	USWs []*USW
	UDMs []*UDM
}

// Config is the data passed into our library. This configures things and allows
// us to connect to a controller and write log messages.
type Config struct {
	User      string
	Pass      string
	URL       string
	VerifySSL bool
	ErrorLog  Logger
	DebugLog  Logger
}

// Unifi is what you get in return for providing a password! Unifi represents
// a controller that you can make authenticated requests to. Use this to make
// additional requests for devices, clients or other custom data. Do not set
// the loggers to nil. Set them to DiscardLogs if you want no logs.
type Unifi struct {
	*http.Client
	*Config
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
	case string:
		f.Txt = i
		f.Val, _ = strconv.ParseFloat(i, 64)
	case nil:
		f.Txt = "0"
		f.Val = 0
	default:
		return fmt.Errorf("cannot unmarshal to FlexInt: %s", b)
	}

	return nil
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
	f.Txt = strings.Trim(string(b), `"`)
	f.Val = f.Txt == "1" || strings.EqualFold(f.Txt, "true") || strings.EqualFold(f.Txt, "yes") ||
		strings.EqualFold(f.Txt, "t") || strings.EqualFold(f.Txt, "armed") || strings.EqualFold(f.Txt, "active") ||
		strings.EqualFold(f.Txt, "enabled") || strings.EqualFold(f.Txt, "ready") || strings.EqualFold(f.Txt, "up") ||
		strings.EqualFold(f.Txt, "ok")

	return nil
}
