package unifipoller

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golift/unifi"
	"github.com/pkg/errors"
)

// DumpJSONPayload prints raw json from the Unifi Controller.
func (u *UnifiPoller) DumpJSONPayload() (err error) {
	u.Quiet = true
	u.Unifi, err = unifi.NewUnifi(u.UnifiUser, u.UnifiPass, u.UnifiBase, u.VerifySSL)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "[INFO] Authenticated to Unifi Controller @", u.UnifiBase, "as user", u.UnifiUser)
	if err := u.CheckSites(); err != nil {
		return err
	}
	u.Unifi.ErrorLog = func(m string, v ...interface{}) {
		fmt.Fprintf(os.Stderr, "[ERROR] "+m, v...)
	} // Log all errors to stderr.

	switch sites, err := u.GetFilteredSites(); {
	case err != nil:
		return err
	case StringInSlice(u.DumpJSON, []string{"d", "device", "devices"}):
		return u.DumpDeviceJSON(sites)
	case StringInSlice(u.DumpJSON, []string{"client", "clients", "c"}):
		return u.DumpClientsJSON(sites)
	case strings.HasPrefix(u.DumpJSON, "other "):
		return u.DumpOtherJSON(sites)
	default:
		return errors.New("must provide filter: devices, clients")
	}
}

// DumpClientsJSON prints the raw json for clients in a Unifi Controller.
func (u *UnifiPoller) DumpClientsJSON(sites []unifi.Site) error {
	for _, s := range sites {
		path := fmt.Sprintf(unifi.ClientPath, s.Name)
		if err := u.dumpJSON(path, "Client", s); err != nil {
			return err
		}
	}
	return nil
}

// DumpDeviceJSON prints the raw json for devices in a Unifi Controller.
func (u *UnifiPoller) DumpDeviceJSON(sites []unifi.Site) error {
	for _, s := range sites {
		path := fmt.Sprintf(unifi.DevicePath, s.Name)
		if err := u.dumpJSON(path, "Device", s); err != nil {
			return err
		}
	}
	return nil
}

// DumpOtherJSON prints the raw json for a user-provided path in a Unifi Controller.
func (u *UnifiPoller) DumpOtherJSON(sites []unifi.Site) error {
	for _, s := range sites {
		path := strings.SplitN(u.DumpJSON, " ", 2)[1]
		if strings.Contains(path, "%s") {
			path = fmt.Sprintf(path, s.Name)
		}
		if err := u.dumpJSON(path, "Other", s); err != nil {
			return err
		}
	}
	return nil
}

func (u *UnifiPoller) dumpJSON(path, what string, site unifi.Site) error {
	req, err := u.UniReq(path, "")
	if err != nil {
		return err
	}
	resp, err := u.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "[INFO] Dumping %s JSON for site %s (%s)\n", what, site.Desc, site.Name)
	fmt.Println(string(body))
	return nil
}
