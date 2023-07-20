package webserver

import (
	"fmt"
	"time"
)

// Logf logs a message.
func (s *Server) Logf(msg string, v ...any) {
	NewOutputEvent(PluginName, PluginName, &Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "info"},
	})

	if s.Collect != nil {
		s.Collect.Logf(msg, v...)
	}
}

// LogErrorf logs an error message.
func (s *Server) LogErrorf(msg string, v ...any) {
	NewOutputEvent(PluginName, PluginName, &Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "error"},
	})

	if s.Collect != nil {
		s.Collect.LogErrorf(msg, v...)
	}
}

// LogDebugf logs a debug message.
func (s *Server) LogDebugf(msg string, v ...any) {
	NewOutputEvent(PluginName, PluginName, &Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "debug"},
	})

	if s.Collect != nil {
		s.Collect.LogDebugf(msg, v...)
	}
}
