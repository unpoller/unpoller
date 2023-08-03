package lokiunifi

import (
	"fmt"
	"time"

	"github.com/unpoller/unpoller/pkg/webserver"
)

// Logf logs a message.
func (l *Loki) Logf(msg string, v ...any) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "info"},
	})
	
	if l.Collect != nil {
		l.Collect.Logf(msg, v...)
	}
}

// LogErrorf logs an error message.
func (l *Loki) LogErrorf(msg string, v ...any) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "error"},
	})
	
	if l.Collect != nil {
		l.Collect.LogErrorf(msg, v...)
	}
}

// LogDebugf logs a debug message.
func (l *Loki) LogDebugf(msg string, v ...any) {
	webserver.NewOutputEvent(PluginName, PluginName, &webserver.Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "debug"},
	})
	
	if l.Collect != nil {
		l.Collect.LogDebugf(msg, v...)
	}
}
