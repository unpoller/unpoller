package poller

import (
	"fmt"
	"log"
	"strings"
)

const callDepth = 2

// LogError logs an error and increments the error counter.
// Should be used in the poller loop.
func (u *UnifiPoller) LogError(err error, prefix string) {
	if err != nil {
		u.errorCount++
		_ = log.Output(callDepth, fmt.Sprintf("[ERROR] %v: %v", prefix, err))
	}
}

// StringInSlice returns true if a string is in a slice.
func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}

// Logf prints a log entry if quiet is false.
func (u *UnifiPoller) Logf(m string, v ...interface{}) {
	if !u.Config.Quiet {
		_ = log.Output(callDepth, fmt.Sprintf("[INFO] "+m, v...))
	}
}

// LogDebugf prints a debug log entry if debug is true and quite is false
func (u *UnifiPoller) LogDebugf(m string, v ...interface{}) {
	if u.Config.Debug && !u.Config.Quiet {
		_ = log.Output(callDepth, fmt.Sprintf("[DEBUG] "+m, v...))
	}
}

// LogErrorf prints an error log entry. This is used for external library logging.
func (u *UnifiPoller) LogErrorf(m string, v ...interface{}) {
	_ = log.Output(callDepth, fmt.Sprintf("[ERROR] "+m, v...))
}
