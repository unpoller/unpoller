package unittest

import (
	"testing"

	"github.com/unpoller/unifi/v5/mocks"
	"github.com/unpoller/unpoller/pkg/inputunifi"
	"github.com/unpoller/unpoller/pkg/poller"
)

type TestRig struct {
	MockServer *mocks.MockHTTPTestServer
	Collector  *poller.TestCollector
	InputUnifi *inputunifi.InputUnifi
	Controller *inputunifi.Controller
}

func NewTestSetup(t *testing.T) *TestRig {
	srv := mocks.NewMockHTTPTestServer()
	testCollector := poller.NewTestCollector(t)

	enabled := true
	controller := inputunifi.Controller{
		SaveAnomal: &enabled,
		SaveAlarms: &enabled,
		SaveEvents: &enabled,
		SaveIDs:    &enabled,
		SaveDPI:    &enabled,
		SaveRogue:  &enabled,
		SaveSites:  &enabled,
		URL:        srv.Server.URL,
	}
	in := &inputunifi.InputUnifi{
		Logger: testCollector.Logger,
		Config: &inputunifi.Config{
			Disable:     false,
			Default:     controller,
			Controllers: []*inputunifi.Controller{&controller},
		},
	}
	testCollector.AddInput(&poller.InputPlugin{
		Name:  "unifi",
		Input: in,
	})

	return &TestRig{
		MockServer: srv,
		Collector:  testCollector,
		InputUnifi: in,
		Controller: &controller,
	}
}

func (t *TestRig) Initialize() {
	_ = t.InputUnifi.Initialize(t.Collector.Logger)
	_, _ = t.InputUnifi.Metrics(nil)
	_, _ = t.InputUnifi.Events(nil)
}

func (t *TestRig) Close() {
	t.MockServer.Server.Close()
}

func PBool(v bool) *bool {
	return &v
}
