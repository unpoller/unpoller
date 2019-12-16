package poller

import (
	"strings"
)

// DumpJSONPayload prints raw json from the UniFi Controller.
// This only works with controller 0 (first one) in the config.
func (u *UnifiPoller) DumpJSONPayload() (err error) {
	if true {
		return nil
	}
	/*
		u.Config.Quiet = true
		config := u.Config.Controllers[0]

		config.Unifi, err = unifi.NewUnifi(&unifi.Config{
			User:      config.User,
			Pass:      config.Pass,
			URL:       config.URL,
			VerifySSL: config.VerifySSL,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "[INFO] Authenticated to UniFi Controller @ %v as user %v", config.URL, config.User)

		if err := u.CheckSites(config); err != nil {
			return err
		}

		config.Unifi.ErrorLog = func(m string, v ...interface{}) {
			fmt.Fprintf(os.Stderr, "[ERROR] "+m, v...)
		} // Log all errors to stderr.

		switch sites, err := u.GetFilteredSites(config); {
		case err != nil:
			return err
		case StringInSlice(u.Flags.DumpJSON, []string{"d", "device", "devices"}):
			return u.dumpSitesJSON(config, unifi.APIDevicePath, "Devices", sites)
		case StringInSlice(u.Flags.DumpJSON, []string{"client", "clients", "c"}):
			return u.dumpSitesJSON(config, unifi.APIClientPath, "Clients", sites)
		case strings.HasPrefix(u.Flags.DumpJSON, "other "):
			apiPath := strings.SplitN(u.Flags.DumpJSON, " ", 2)[1]
			_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping Path '%s':\n", apiPath)
			return u.PrintRawAPIJSON(config, apiPath)
		default:
			return fmt.Errorf("must provide filter: devices, clients, other")
		}
	*/
	return nil
}

/*
func (u *UnifiPoller) dumpSitesJSON(c Controller, path, name string, sites unifi.Sites) error {
	for _, s := range sites {
		apiPath := fmt.Sprintf(path, s.Name)
		_, _ = fmt.Fprintf(os.Stderr, "[INFO] Dumping %s: '%s' JSON for site: %s (%s):\n",
			name, apiPath, s.Desc, s.Name)
		if err := u.PrintRawAPIJSON(c, apiPath); err != nil {
			return err
		}
	}
	return nil
}

// PrintRawAPIJSON prints the raw json for a user-provided path on a UniFi Controller.
func (u *UnifiPoller) PrintRawAPIJSON(c Controller, apiPath string) error {
	body, err := c.Unifi.GetJSON(apiPath)
	fmt.Println(string(body))
	return err
}
*/

// StringInSlice returns true if a string is in a slice.
func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}

	return false
}
