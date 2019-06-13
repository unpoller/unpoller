package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golift/unifi"
	"github.com/pkg/errors"
)

// DumpJSON prints raw json from the Unifi Controller.
func (c *Config) DumpJSON(filter string) error {
	c.Quiet = true
	controller, err := unifi.NewUnifi(c.UnifiUser, c.UnifiPass, c.UnifiBase, c.VerifySSL)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Authenticated to Unifi Controller @", c.UnifiBase, "as user", c.UnifiUser)
	if err := c.CheckSites(controller); err != nil {
		return err
	}
	controller.ErrorLog = func(m string, v ...interface{}) {
		fmt.Fprintf(os.Stderr, "[ERROR] "+m, v...)
	} // Log all errors to stderr.

	switch sites, err := filterSites(controller, c.Sites); {
	case err != nil:
		return err
	case StringInSlice(filter, []string{"d", "device", "devices"}):
		return c.DumpDeviceJSON(sites, controller)
	case StringInSlice(filter, []string{"client", "clients", "c"}):
		return c.DumpClientsJSON(sites, controller)
	default:
		return errors.New("must provide filter: devices, clients")
	}
}

// DumpClientsJSON prints the raw json for clients in a Unifi Controller.
func (c *Config) DumpClientsJSON(sites []unifi.Site, controller *unifi.Unifi) error {
	for _, s := range sites {
		path := fmt.Sprintf(unifi.ClientPath, s.Name)
		if err := dumpJSON(path, "Client", s, controller); err != nil {
			return err
		}
	}
	return nil
}

// DumpDeviceJSON prints the raw json for devices in a Unifi Controller.
func (c *Config) DumpDeviceJSON(sites []unifi.Site, controller *unifi.Unifi) error {
	for _, s := range sites {
		path := fmt.Sprintf(unifi.DevicePath, s.Name)
		if err := dumpJSON(path, "Device", s, controller); err != nil {
			return err
		}
	}
	return nil
}

func dumpJSON(path, what string, site unifi.Site, controller *unifi.Unifi) error {
	req, err := controller.UniReq(path, "")
	if err != nil {
		return err
	}
	resp, err := controller.Do(req)
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
	fmt.Fprintf(os.Stderr, "Dumping %s JSON for site %s (%s)\n", what, site.Desc, site.Name)
	fmt.Println(string(body))
	return nil
}
