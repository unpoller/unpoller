package promunifi

import (
	"fmt"
	"time"

	"github.com/unpoller/webserver"
)

// Logf logs a message.
func (u *promUnifi) Logf(msg string, v ...interface{}) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "info"},
	})
	u.Collector.Logf(msg, v...)
}

// LogErrorf logs an error message.
func (u *promUnifi) LogErrorf(msg string, v ...interface{}) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "error"},
	})
	u.Collector.LogErrorf(msg, v...)
}

// LogDebugf logs a debug message.
func (u *promUnifi) LogDebugf(msg string, v ...interface{}) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "debug"},
	})
	u.Collector.LogDebugf(msg, v...)
}
