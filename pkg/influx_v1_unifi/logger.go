package influx_v1_unifi

import (
	"fmt"
	"time"

	"github.com/unpoller/unpoller/pkg/webserver"
)

// Logf logs a message.
func (u *InfluxV1Unifi) Logf(msg string, v ...any) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "info"},
	})
	u.Collector.Logf(msg, v...)
}

// LogErrorf logs an error message.
func (u *InfluxV1Unifi) LogErrorf(msg string, v ...any) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "error"},
	})
	u.Collector.LogErrorf(msg, v...)
}

// LogDebugf logs a debug message.
func (u *InfluxV1Unifi) LogDebugf(msg string, v ...any) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "debug"},
	})
	u.Collector.LogDebugf(msg, v...)
}
