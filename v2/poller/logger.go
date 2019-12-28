package poller

import (
	"fmt"
	"log"
)

const callDepth = 2

// Logger is passed into input packages so they may write logs.
type Logger interface {
	Logf(m string, v ...interface{})
	LogErrorf(m string, v ...interface{})
	LogDebugf(m string, v ...interface{})
}

// Logf prints a log entry if quiet is false.
func (u *UnifiPoller) Logf(m string, v ...interface{}) {
	if !u.Quiet {
		_ = log.Output(callDepth, fmt.Sprintf("[INFO] "+m, v...))
	}
}

// LogDebugf prints a debug log entry if debug is true and quite is false
func (u *UnifiPoller) LogDebugf(m string, v ...interface{}) {
	if u.Debug && !u.Quiet {
		_ = log.Output(callDepth, fmt.Sprintf("[DEBUG] "+m, v...))
	}
}

// LogErrorf prints an error log entry.
func (u *UnifiPoller) LogErrorf(m string, v ...interface{}) {
	_ = log.Output(callDepth, fmt.Sprintf("[ERROR] "+m, v...))
}
