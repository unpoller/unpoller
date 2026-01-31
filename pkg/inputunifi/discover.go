package inputunifi

import (
	"fmt"
)

// Discover implements poller.Discoverer. It uses the first configured controller
// to probe known API endpoints and write a shareable report to outputPath.
// Uses the same credentials as normal polling (from config file).
func (u *InputUnifi) Discover(outputPath string) error {
	if u.Config == nil || u.Disable {
		return fmt.Errorf("unifi input disabled or not configured")
	}

	u.setDefaults(&u.Default)

	if len(u.Controllers) == 0 && !u.Dynamic {
		u.Controllers = []*Controller{&u.Default}
	}

	if len(u.Controllers) == 0 {
		return fmt.Errorf("no unifi controller configured")
	}

	c := u.setControllerDefaults(u.Controllers[0])
	if c.URL == "" {
		return fmt.Errorf("first controller has no URL")
	}

	if err := u.getUnifi(c); err != nil {
		return fmt.Errorf("authenticating to controller: %w", err)
	}

	sites, err := c.Unifi.GetSites()
	if err != nil {
		return fmt.Errorf("getting sites: %w", err)
	}

	site := "default"
	if len(sites) > 0 && sites[0].Name != "" {
		site = sites[0].Name
	}

	if err := c.Unifi.DiscoverEndpoints(site, outputPath); err != nil {
		return fmt.Errorf("writing discovery report: %w", err)
	}

	return nil
}
