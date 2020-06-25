package lokiunifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	lokiPushPath = "/loki/api/v1/push"
)

var (
	errStatusCode = fmt.Errorf("unexpected HTTP status code")
)

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
func (c *Client) Post(logs LogStreams) error {
	msg, err := json.Marshal(logs)
	if err != nil {
		return err
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

		return errors.Wrap(errStatusCode, m)
	}

	return nil
}

// NewRequest creates the http request based on input data.
func (c *Client) NewRequest(url, method, cType string, msg []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(msg))
	if err != nil {
		return nil, err
	}

	if cType != "" {
		req.Header.Set("Content-Type", cType)
	}

	if c.Config.Username != "" || c.Config.Password != "" {
		req.SetBasicAuth(c.Config.Username, c.Config.Password)
	}

	if c.Config.TenantID != "" {
		req.Header.Set("X-Scope-OrgID", c.Config.TenantID)
	}

	return req, nil
}

// Do makes an http request and returns the status code, body and/or an error.
func (c *Client) Do(req *http.Request) (int, []byte, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, body, err
}
