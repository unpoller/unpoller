package unifi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlexInt(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	type testReply struct {
		Five    FlexInt `json:"five"`
		Seven   FlexInt `json:"seven"`
		Auto    FlexInt `json:"auto"`
		Channel FlexInt `json:"channel"`
	}
	var r testReply
	// test unmarshalling the custom type three times with different values.
	a.Nil(json.Unmarshal([]byte(`{"five": "5", "seven": 7, "auto": "auto"}`), &r))

	// test number in string.
	a.EqualValues(5, r.Five.Number)
	a.EqualValues("5", r.Five.String)
	// test number.
	a.EqualValues(7, r.Seven.Number)
	a.EqualValues("7", r.Seven.String)
	// test string.
	a.EqualValues(0, r.Auto.Number)
	a.EqualValues("auto", r.Auto.String)
	// test (error) struct.
	a.NotNil(json.Unmarshal([]byte(`{"channel": {}}`), &r),
		"a non-string and non-number must produce an error.")
	a.EqualValues(0, r.Channel.Number)
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

func TestAuthController(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	url := "http://127.0.0.1:64431"
	authReq, err := AuthController("user1", "pass2", url, false)
	a.NotNil(err)
	a.EqualValues(url, authReq.baseURL)
	a.Contains(err.Error(), "authReq.Do(req):", "an invalid destination should product a .Do(req) error.")
	/* TODO: OPEN web server, check parameters posted, more. This test is incomplete.
	a.EqualValues(`{"username": "user1","password": "pass2"}`, string(post_params), "user/pass json parameters improperly encoded")
	*/
}
