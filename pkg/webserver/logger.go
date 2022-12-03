package webserver

import (
	"fmt"
	"time"
)

// Logf logs a message.
func (s *Server) Logf(msg string, v ...interface{}) {
	NewOutputEvent(PluginName, PluginName, &Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "info"},
	})
	s.Collect.Logf(msg, v...)
}

// LogErrorf logs an error message.
func (s *Server) LogErrorf(msg string, v ...interface{}) {
	NewOutputEvent(PluginName, PluginName, &Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "error"},
	})
	s.Collect.LogErrorf(msg, v...)
}

// LogDebugf logs a debug message.
func (s *Server) LogDebugf(msg string, v ...interface{}) {
	NewOutputEvent(PluginName, PluginName, &Event{
		Ts:   time.Now(),
		Msg:  fmt.Sprintf(msg, v...),
		Tags: map[string]string{"type": "debug"},
	})
	s.Collect.LogDebugf(msg, v...)
}
