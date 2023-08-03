package poller

import (
	"sync"
)

type TestCollector struct {
	sync.RWMutex
	Logger Logger
	inputs []*InputPlugin
	poller *Poller
}

func NewTestCollector(l testLogger) *TestCollector {
	return &TestCollector{
		Logger: NewTestLogger(l),
	}
}

func (t *TestCollector) AddInput(input *InputPlugin) {
	t.Lock()
	defer t.Unlock()

	t.inputs = append(t.inputs, input)
}

func (t *TestCollector) Metrics(filter *Filter) (*Metrics, error) {
	t.RLock()
	defer t.RUnlock()

	return collectMetrics(filter, t.inputs)
}

func (t *TestCollector) Events(filter *Filter) (*Events, error) {
	t.RLock()
	defer t.RUnlock()

	return collectEvents(filter, t.inputs)
}

func (t *TestCollector) SetPoller(poller *Poller) {
	t.Lock()
	defer t.Unlock()

	t.poller = poller
}

func (t *TestCollector) Poller() Poller {
	return *t.poller
}

func (t *TestCollector) Inputs() (names []string) {
	t.RLock()
	defer t.RUnlock()

	for i := range t.inputs {
		names = append(names, inputs[i].Name)
	}

	return names
}

func (t *TestCollector) Outputs() []string {
	return []string{}
}

func (t *TestCollector) Logf(m string, v ...any) {
	t.Logger.Logf(m, v...)
}

func (t *TestCollector) LogErrorf(m string, v ...any) {
	t.Logger.LogErrorf(m, v...)
}

func (t *TestCollector) LogDebugf(m string, v ...any) {
	t.Logger.LogDebugf(m, v...)
}


type testLogger interface {
	Log(args ...any)
	Logf(format string, args ...any)
}

type TestLogger struct {
	log testLogger
}

func NewTestLogger(l testLogger) *TestLogger {
	return &TestLogger{log: l}
}

// Logf prints a log entry if quiet is false.
func (t *TestLogger) Logf(m string, v ...any) {
	t.log.Logf("[INFO] "+m, v...)
}

// LogDebugf prints a debug log entry if debug is true and quite is false.
func (t *TestLogger) LogDebugf(m string, v ...any) {
	t.log.Logf("[DEBUG] "+m, v...)
}

// LogErrorf prints an error log entry.
func (t *TestLogger) LogErrorf(m string, v ...any) {
	t.log.Logf("[ERROR] "+m, v...)
}
