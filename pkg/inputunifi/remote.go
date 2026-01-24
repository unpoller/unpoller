package inputunifi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	remoteAPIBaseURL = "https://api.ui.com"
	remoteAPIVersion = "v1"
)

// Console represents a UniFi console from the remote API.
type Console struct {
	ID          string `json:"id"`
	IPAddress   string `json:"ipAddress"`
	Type        string `json:"type"`
	Owner       bool   `json:"owner"`
	IsBlocked   bool   `json:"isBlocked"`
	ReportedState struct {
		Name      string `json:"name"`
		Hostname  string `json:"hostname"`
		IP        string `json:"ip"`
		State     string `json:"state"`
		Mac       string `json:"mac"`
	} `json:"reportedState"`
	ConsoleName string // Derived field: name from reportedState
}

// HostsResponse represents the response from /v1/hosts endpoint.
type HostsResponse struct {
	Data           []Console `json:"data"`
	HTTPStatusCode int       `json:"httpStatusCode"`
	TraceID        string    `json:"traceId"`
	NextToken      string    `json:"nextToken,omitempty"`
}

// Site represents a site from the remote API.
type RemoteSite struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SitesResponse represents the response from the sites endpoint.
type SitesResponse struct {
	Data           []RemoteSite `json:"data"`
	HTTPStatusCode int          `json:"httpStatusCode"`
	TraceID        string       `json:"traceId"`
}

// remoteAPIClient handles HTTP requests to the remote UniFi API.
type remoteAPIClient struct {
	apiKey   string
	baseURL  string
	client   *http.Client
	logError func(string, ...any)
	logDebug func(string, ...any)
	log      func(string, ...any)
}

// newRemoteAPIClient creates a new remote API client.
func (u *InputUnifi) newRemoteAPIClient(apiKey string) *remoteAPIClient {
	if apiKey == "" {
		return nil
	}

	// Handle file:// prefix for API key
	if strings.HasPrefix(apiKey, "file://") {
		apiKey = u.getPassFromFile(strings.TrimPrefix(apiKey, "file://"))
	}

	return &remoteAPIClient{
		apiKey:  apiKey,
		baseURL: remoteAPIBaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
		logError: u.LogErrorf,
		logDebug: u.LogDebugf,
		log:      u.Logf,
	}
}

// makeRequest makes an HTTP request to the remote API.
func (c *remoteAPIClient) makeRequest(method, path string, queryParams map[string]string) ([]byte, error) {
	fullURL := c.baseURL + path

	if len(queryParams) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, fmt.Errorf("parsing URL: %w", err)
		}

		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	c.logDebug("Making %s request to: %s", method, fullURL)

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// discoverConsoles discovers all consoles available via the remote API.
func (c *remoteAPIClient) discoverConsoles() ([]Console, error) {
	// Start with first page
	queryParams := map[string]string{
		"pageSize": "10",
	}

	var allConsoles []Console
	nextToken := ""

	for {
		if nextToken != "" {
			queryParams["nextToken"] = nextToken
		} else {
			// Remove nextToken from params for first request
			delete(queryParams, "nextToken")
		}

		body, err := c.makeRequest("GET", "/v1/hosts", queryParams)
		if err != nil {
			return nil, fmt.Errorf("fetching consoles: %w", err)
		}

		var response HostsResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("parsing consoles response: %w", err)
		}

		// Filter for console type only
		for _, console := range response.Data {
			if console.Type == "console" && !console.IsBlocked {
				// Extract the console name from reportedState
				console.ConsoleName = console.ReportedState.Name
				if console.ConsoleName == "" {
					console.ConsoleName = console.ReportedState.Hostname
				}
				allConsoles = append(allConsoles, console)
			}
		}

		// Check if there's a nextToken to continue pagination
		if response.NextToken == "" {
			break
		}

		nextToken = response.NextToken
		c.logDebug("Fetching next page of consoles with nextToken: %s", nextToken)
	}

	return allConsoles, nil
}

// discoverSites discovers all sites for a given console ID.
func (c *remoteAPIClient) discoverSites(consoleID string) ([]RemoteSite, error) {
	path := fmt.Sprintf("/v1/connector/consoles/%s/proxy/network/integration/v1/sites", consoleID)

	queryParams := map[string]string{
		"offset": "0",
		"limit":  "100",
	}

	body, err := c.makeRequest("GET", path, queryParams)
	if err != nil {
		return nil, fmt.Errorf("fetching sites for console %s: %w", consoleID, err)
	}

	var response SitesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parsing sites response: %w", err)
	}

	return response.Data, nil
}

// discoverRemoteControllers discovers all controllers via remote API and creates Controller entries.
func (u *InputUnifi) discoverRemoteControllers(apiKey string) ([]*Controller, error) {
	client := u.newRemoteAPIClient(apiKey)
	if client == nil {
		return nil, fmt.Errorf("remote API key not provided")
	}

	u.Logf("Discovering remote UniFi consoles...")

	consoles, err := client.discoverConsoles()
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

		sites, err := client.discoverSites(console.ID)
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
			URL: fmt.Sprintf("%s/v1/connector/consoles/%s", remoteAPIBaseURL, console.ID),
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
		if len(siteNames) == 1 && strings.EqualFold(siteNames[0], "default") && consoleName != "" {
			controller.DefaultSiteNameOverride = consoleName
			// Set sites to "all" since we're overriding the default site name
			controller.Sites = []string{"all"}
			u.LogDebugf("Using console name '%s' as default site name override for Cloud Gateway", consoleName)
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
