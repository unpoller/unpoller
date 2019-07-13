package unifipoller

import (
	"fmt"
	"os"
	"strings"

	"github.com/golift/unifi"
	"github.com/pkg/errors"
)

// DumpJSONPayload prints raw json from the UniFi Controller.
func (u *UnifiPoller) DumpJSONPayload() (err error) {
	u.Quiet = true
	u.Unifi, err = unifi.NewUnifi(u.UnifiUser, u.UnifiPass, u.UnifiBase, u.VerifySSL)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "[INFO] Authenticated to UniFi Controller @", u.UnifiBase, "as user", u.UnifiUser)
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
		return u.dumpSitesJSON(unifi.DevicePath, "Devices", sites)
	case StringInSlice(u.DumpJSON, []string{"client", "clients", "c"}):
		return u.dumpSitesJSON(unifi.ClientPath, "Clients", sites)
	case strings.HasPrefix(u.DumpJSON, "other "):
		apiPath := strings.SplitN(u.DumpJSON, " ", 2)[1]
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping Path '%s':\n", apiPath)
		return u.PrintRawAPIJSON(apiPath)
	default:
		return errors.New("must provide filter: devices, clients, other")
	}
}

func (u *UnifiPoller) dumpSitesJSON(path, name string, sites unifi.Sites) error {
	for _, s := range sites {
		apiPath := fmt.Sprintf(path, s.Name)
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping %s: '%s' JSON for site: %s (%s):\n", name, apiPath, s.Desc, s.Name)
		if err := u.PrintRawAPIJSON(apiPath); err != nil {
			return err
		}
	}
	return nil
}

// PrintRawAPIJSON prints the raw json for a user-provided path on a UniFi Controller.
func (u *UnifiPoller) PrintRawAPIJSON(apiPath string) error {
	body, err := u.GetJSON(apiPath)
	fmt.Println(string(body))
	return err
}
