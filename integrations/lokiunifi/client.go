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

type client struct {
	*Config
	*http.Client
}

func (l *Loki) httpClient() *client {
	return &client{
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

// Send marshals and posts a batch of log messages.
func (c *client) Send(logs Logs) error {
	msg, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	u := strings.TrimSuffix(c.URL, lokiPushPath) + lokiPushPath

	code, body, err := c.PostReq(u, "application/json", msg)
	if err != nil {
		return err
	} else if code != http.StatusNoContent {
		m := fmt.Sprintf("%s (%d/%s) %s, msg: %s", u, code, http.StatusText(code),
			strings.TrimSpace(strings.ReplaceAll(string(body), "\n", " ")), msg)

		return errors.Wrap(errStatusCode, m)
	}

	return nil
}

// PostReq posts data to a url with a custom content type.
// Returns the status code, body and/or an error.
func (c *client) PostReq(url, cType string, msg []byte) (int, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", cType)

	resp, err := c.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, body, err
}
