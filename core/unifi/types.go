package unifi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrCannotUnmarshalFlexInt = fmt.Errorf("cannot unmarshal to FlexInt")

// This is a list of unifi API paths.
// The %s in each string must be replaced with a Site.Name.
const (
	// APIRogueAP shows your neighbors' wifis.
	APIRogueAP string = "/api/s/%s/stat/rogueap"
	// APIStatusPath shows Controller version.
	APIStatusPath string = "/status"
	// APIEventPath contains UniFi Event data.
	APIEventPath string = "/api/s/%s/stat/event"
	// APISiteList is the path to the api site list.
	APISiteList string = "/api/stat/sites"
	// APISiteDPI is site DPI data.
	APISiteDPI string = "/api/s/%s/stat/sitedpi"
	// APISiteDPI is site DPI data.
	APIClientDPI string = "/api/s/%s/stat/stadpi"
	// APIClientPath is Unifi Clients API Path.
	APIClientPath string = "/api/s/%s/stat/sta"
	// APIAllUserPath is Unifi Insight all previous Clients API Path.
	APIAllUserPath string = "/api/s/%s/stat/alluser"
	// APINetworkPath is where we get data about Unifi networks.
	APINetworkPath string = "/api/s/%s/rest/networkconf"
	// APIDevicePath is where we get data about Unifi devices.
	APIDevicePath string = "/api/s/%s/stat/device"
	// APILoginPath is Unifi Controller Login API Path.
	APILoginPath string = "/api/login"
	// APILoginPathNew is how we log into UDM 5.12.55+.
	APILoginPathNew string = "/api/auth/login"
	// APILogoutPath is how we logout from UDM.
	APILogoutPath string = "/api/logout"
	// APIEventPathIDS returns Intrusion Detection/Prevention Systems Events.
	APIEventPathIDS string = "/api/s/%s/stat/ips/event"
	// APIEventPathAlarms contains the site alarms.
	APIEventPathAlarms string = "/api/s/%s/list/alarm"
	// APIPrefixNew is the prefix added to the new API paths; except login. duh.
	APIPrefixNew string = "/proxy/network"
	// APIAnomaliesPath returns site anomalies.
	APIAnomaliesPath string = "/api/s/%s/stat/anomalies"
	APICommandPath   string = "/api/s/%s/cmd"
	APIDevMgrPath    string = APICommandPath + "/devmgr"
)

// path returns the correct api path based on the new variable.
// new is based on the unifi-controller output. is it new or old output?
func (u *Unifi) path(path string) string {
	if u.new {
		if path == APILoginPath {
			return APILoginPathNew
		}

		if !strings.HasPrefix(path, APIPrefixNew) && path != APILoginPathNew {
			return APIPrefixNew + path
		}
	}

	return path
}

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
	UXGs []*UXG
}

// Config is the data passed into our library. This configures things and allows
// us to connect to a controller and write log messages. Optional SSLCert is used
// for ssl cert pinning; provide the content of a PEM to validate the server's cert.
type Config struct {
	User      string
	Pass      string
	URL       string
	SSLCert   [][]byte
	ErrorLog  Logger
	DebugLog  Logger
	Timeout   time.Duration // how long to wait for replies, default: forever.
	VerifySSL bool
}

// Unifi is what you get in return for providing a password! Unifi represents
// a controller that you can make authenticated requests to. Use this to make
// additional requests for devices, clients or other custom data. Do not set
// the loggers to nil. Set them to DiscardLogs if you want no logs.
type Unifi struct {
	*http.Client
	*Config
	*server
	csrf         string
	fingerprints fingerprints
	new          bool
}

type fingerprints []string

// Contains returns true if the fingerprint is in the list.
func (f fingerprints) Contains(s string) bool {
	for i := range f {
		if s == f[i] {
			return true
		}
	}

	return false
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
		return fmt.Errorf("json unmarshal: %w", err)
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
		return fmt.Errorf("%v: %w", b, ErrCannotUnmarshalFlexInt)
	}

	return nil
}

func (f *FlexInt) String() string {
	return f.Txt
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

func (f *FlexBool) String() string {
	return f.Txt
}
