package poller

import (
	"fmt"
	"os"
	"strings"

	"golift.io/unifi"
)

// DumpJSONPayload prints raw json from the UniFi Controller.
func (u *UnifiPoller) DumpJSONPayload() (err error) {
	u.Config.Quiet = true
	u.Unifi, err = unifi.NewUnifi(&unifi.Config{
		User:      u.Config.UnifiUser,
		Pass:      u.Config.UnifiPass,
		URL:       u.Config.UnifiBase,
		VerifySSL: u.Config.VerifySSL,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "[INFO] Authenticated to UniFi Controller @ %v as user %v",
		u.Config.UnifiBase, u.Config.UnifiUser)
	if err := u.CheckSites(); err != nil {
		return err
	}

	u.Unifi.ErrorLog = func(m string, v ...interface{}) {
		fmt.Fprintf(os.Stderr, "[ERROR] "+m, v...)
	} // Log all errors to stderr.

	switch sites, err := u.GetFilteredSites(); {
	case err != nil:
		return err
	case StringInSlice(u.Flag.DumpJSON, []string{"d", "device", "devices"}):
		return u.dumpSitesJSON(unifi.APIDevicePath, "Devices", sites)
	case StringInSlice(u.Flag.DumpJSON, []string{"client", "clients", "c"}):
		return u.dumpSitesJSON(unifi.APIClientPath, "Clients", sites)
	case strings.HasPrefix(u.Flag.DumpJSON, "other "):
		apiPath := strings.SplitN(u.Flag.DumpJSON, " ", 2)[1]
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping Path '%s':\n", apiPath)
		return u.PrintRawAPIJSON(apiPath)
	default:
		return fmt.Errorf("must provide filter: devices, clients, other")
	}
}

func (u *UnifiPoller) dumpSitesJSON(path, name string, sites unifi.Sites) error {
	for _, s := range sites {
		apiPath := fmt.Sprintf(path, s.Name)
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping %s: '%s' JSON for site: %s (%s):\n",
			name, apiPath, s.Desc, s.Name)
		if err := u.PrintRawAPIJSON(apiPath); err != nil {
			return err
		}
	}
	return nil
}

// PrintRawAPIJSON prints the raw json for a user-provided path on a UniFi Controller.
func (u *UnifiPoller) PrintRawAPIJSON(apiPath string) error {
	body, err := u.Unifi.GetJSON(apiPath)
	fmt.Println(string(body))
	return err
}
