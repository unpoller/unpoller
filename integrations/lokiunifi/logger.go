package lokiunifi

import (
	"fmt"
	"time"

	"github.com/unpoller/unpoller/core/webserver"
)

// Logf logs a message.
func (l *Loki) Logf(msg string, v ...interface{}) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "info"},
	})
	l.Collect.Logf(msg, v...)
}

// LogErrorf logs an error message.
func (l *Loki) LogErrorf(msg string, v ...interface{}) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "error"},
	})
	l.Collect.LogErrorf(msg, v...)
}

// LogDebugf logs a debug message.
func (l *Loki) LogDebugf(msg string, v ...interface{}) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "debug"},
	})
	l.Collect.LogDebugf(msg, v...)
}
