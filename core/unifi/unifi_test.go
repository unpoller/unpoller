package unifi

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUnifi(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	url := "http://127.0.0.1:64431"
	authReq, err := NewUnifi("user1", "pass2", url, false)
	a.NotNil(err)
	a.EqualValues(url, authReq.baseURL)
	a.Contains(err.Error(), "authReq.Do(req):", "an invalid destination should product a .Do(req) error.")
	/* TODO: OPEN web server, check parameters posted, more. This test is incomplete.
	a.EqualValues(`{"username": "user1","password": "pass2"}`, string(post_params),
		"user/pass json parameters improperly encoded")
	*/
}

func TestUniReq(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	u := "/test/path"
	url := "http://some.url:8443"
	// Test empty parameters.
	authReq := &Unifi{Client: &http.Client{}, baseURL: url}
	r, err := authReq.UniReq(u, "")
	a.Nil(err, "newrequest must not produce an error")
	a.EqualValues(u, r.URL.Path,
		"the provided apiPath was not added to http request")
	a.EqualValues(url, r.URL.Scheme+"://"+r.URL.Host, "URL improperly encoded")
	a.EqualValues("GET", r.Method, "without parameters the method must be GET")
	a.EqualValues("application/json", r.Header.Get("Accept"), "Accept header must be set to application/json")

	// Test with parameters
	p := "key1=value9&key2=value7"
	authReq = &Unifi{Client: &http.Client{}, baseURL: "http://some.url:8443"}
	r, err = authReq.UniReq(u, p)
	a.Nil(err, "newrequest must not produce an error")
	a.EqualValues(u, r.URL.Path,
		"the provided apiPath was not added to http request")
	a.EqualValues(url, r.URL.Scheme+"://"+r.URL.Host, "URL improperly encoded")
	a.EqualValues("POST", r.Method, "with parameters the method must be POST")
	a.EqualValues("application/json", r.Header.Get("Accept"), "Accept header must be set to application/json")
	// Check the parameters.
	d, err := ioutil.ReadAll(r.Body)
	a.Nil(err, "problem reading request body, POST parameters may be malformed")
	a.EqualValues(p, string(d), "POST parameters improperly encoded")
}
