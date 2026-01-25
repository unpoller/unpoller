package inputunifi

import (
	"fmt"
	"strings"

	"github.com/unpoller/unifi/v5"
)

// discoverRemoteControllers discovers all controllers via remote API and creates Controller entries.
func (u *InputUnifi) discoverRemoteControllers(apiKey string) ([]*Controller, error) {
	// Handle file:// prefix for API key
	if strings.HasPrefix(apiKey, "file://") {
		apiKey = u.getPassFromFile(strings.TrimPrefix(apiKey, "file://"))
	}

	if apiKey == "" {
		return nil, fmt.Errorf("remote API key not provided")
	}

	// Use library client
	client := unifi.NewRemoteAPIClient(apiKey, u.LogErrorf, u.LogDebugf, u.Logf)

	u.Logf("Discovering remote UniFi consoles...")

	consoles, err := client.DiscoverConsoles()
	if err != nil {
		return nil, fmt.Errorf("discovering consoles: %w", err)
	}

	if len(consoles) == 0 {
		u.Logf("No consoles found via remote API")
		return nil, nil
	}

	u.Logf("Found %d console(s) via remote API", len(consoles))

	var controllers []*Controller

	for _, console := range consoles {
		consoleName := console.ConsoleName
		if consoleName == "" {
			consoleName = console.ReportedState.Name
		}
		if consoleName == "" {
			consoleName = console.ReportedState.Hostname
		}
		u.LogDebugf("Discovering sites for console: %s (%s)", console.ID, consoleName)

		sites, err := client.DiscoverSites(console.ID)
		if err != nil {
			u.LogErrorf("Failed to discover sites for console %s: %v", console.ID, err)
			continue
		}

		if len(sites) == 0 {
			u.LogDebugf("No sites found for console %s", console.ID)
			continue
		}

		// Create a controller entry for this console
		// For remote API, the base URL should point to the connector endpoint
		// The unifi library will append /proxy/network/... paths, so we need to account for that
		// For remote API with integration endpoints, we set it to the connector base
		controller := &Controller{
			Remote:    true,
			ConsoleID: console.ID,
			APIKey:    apiKey,
			// Set URL to connector base - the library appends /proxy/network/status
			// But for integration API we need /proxy/network/integration/v1/...
			// This may require library updates, but try connector base first
			URL: fmt.Sprintf("%s/v1/connector/consoles/%s", unifi.RemoteAPIBaseURL, console.ID),
		}

		// Ensure defaults are set before calling setControllerDefaults
		u.setDefaults(&u.Default)

		// Copy defaults
		controller = u.setControllerDefaults(controller)

		// Set remote-specific defaults and ensure all boolean pointers are initialized
		t := true
		f := false
		if controller.VerifySSL == nil {
			controller.VerifySSL = &t // Remote API should verify SSL
		}
		// Ensure all boolean pointers are set (safety check)
		if controller.HashPII == nil {
			controller.HashPII = &f
		}
		if controller.DropPII == nil {
			controller.DropPII = &f
		}
		if controller.SaveSites == nil {
			controller.SaveSites = &t
		}
		if controller.SaveDPI == nil {
			controller.SaveDPI = &f
		}
		if controller.SaveEvents == nil {
			controller.SaveEvents = &f
		}
		if controller.SaveAlarms == nil {
			controller.SaveAlarms = &f
		}
		if controller.SaveAnomal == nil {
			controller.SaveAnomal = &f
		}
		if controller.SaveIDs == nil {
			controller.SaveIDs = &f
		}
		if controller.SaveTraffic == nil {
			controller.SaveTraffic = &f
		}
		if controller.SaveRogue == nil {
			controller.SaveRogue = &f
		}
		if controller.SaveSyslog == nil {
			controller.SaveSyslog = &f
		}
		if controller.SaveProtectLogs == nil {
			controller.SaveProtectLogs = &f
		}
		if controller.ProtectThumbnails == nil {
			controller.ProtectThumbnails = &f
		}

		// Extract site names
		siteNames := make([]string, 0, len(sites))
		for _, site := range sites {
			if site.Name != "" {
				siteNames = append(siteNames, site.Name)
			}
		}

		// For Cloud Gateways, if the only site is "default", use the console name from hosts response
		// as the default site name override. The console name is in reportedState.name
		// (consoleName was already set above in the loop)

		// If we only have one site and it's "default" (case-insensitive), use the console name as override
		// Note: We keep the actual site name ("default") for API calls, but set the override
		// for display/metric naming purposes.
		if len(siteNames) == 1 && strings.EqualFold(siteNames[0], "default") && consoleName != "" {
			controller.DefaultSiteNameOverride = consoleName
			// Keep the actual site name for API calls
			controller.Sites = siteNames
			u.LogDebugf("Using console name '%s' as default site name override for Cloud Gateway (API will use 'default')", consoleName)
		} else if len(siteNames) > 0 {
			controller.Sites = siteNames
		} else {
			controller.Sites = []string{"all"}
		}

		controller.ID = console.ID
		controllers = append(controllers, controller)

		u.Logf("Discovered console %s with %d site(s): %v", consoleName, len(sites), siteNames)
	}

	return controllers, nil
}
