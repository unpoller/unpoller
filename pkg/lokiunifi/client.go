package lokiunifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	lokiPushPath = "/loki/api/v1/push"
)

var errStatusCode = fmt.Errorf("unexpected HTTP status code")

// Client holds the http client for contacting Loki.
type Client struct {
	*Config
	*http.Client
}

func (l *Loki) httpClient() *Client {
	return &Client{
		Config: l.Config,
		Client: &http.Client{
			Timeout: l.Timeout.Duration,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: !l.VerifySSL, // nolint: gosec
				},
			},
		},
	}
}

// Post marshals and posts a batch of log messages.
func (c *Client) Post(logs any) error {
	msg, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	u := strings.TrimSuffix(c.URL, lokiPushPath) + lokiPushPath

	req, err := c.NewRequest(u, "POST", "application/json", msg)
	if err != nil {
		return err
	}

	if code, body, err := c.Do(req); err != nil {
		return err
	} else if code != http.StatusNoContent {
		m := fmt.Sprintf("%s (%d/%s) %s, msg: %s", u, code, http.StatusText(code),
			strings.TrimSpace(strings.ReplaceAll(string(body), "\n", " ")), msg)

		return fmt.Errorf("%s: %w", m, errStatusCode)
	}

	return nil
}

// NewRequest creates the http request based on input data.
func (c *Client) NewRequest(url, method, cType string, msg []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(msg)) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if cType != "" {
		req.Header.Set("Content-Type", cType)
	}

	if c.Username != "" || c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	if c.TenantID != "" {
		req.Header.Set("X-Scope-OrgID", c.TenantID)
	}

	return req, nil
}

// Do makes an http request and returns the status code, body and/or an error.
func (c *Client) Do(req *http.Request) (int, []byte, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, body, fmt.Errorf("reading body: %w", err)
	}

	return resp.StatusCode, body, nil
}
