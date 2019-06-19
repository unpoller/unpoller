package unifipoller

import (
	"fmt"
	"log"
	"strings"
)

// hasErr checks a list of errors for a non-nil.
func hasErr(errs []error) bool {
	for _, err := range errs {
		if err != nil {
			return true
		}
	}
	return false
}

// LogErrors writes a slice of errors, with a prefix, to log-out.
// It also increments the error counter.
func (u *UnifiPoller) LogErrors(errs []error, prefix string) {
	for _, err := range errs {
		if err != nil {
			u.errorCount++
			_ = log.Output(2, fmt.Sprintf("[ERROR] (%v/%v) %v: %v", u.errorCount, u.MaxErrors, prefix, err))
		}
	}
}

// StringInSlice returns true if a string is in a slice.
func StringInSlice(str string, slc []string) bool {
	for _, s := range slc {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}

// Logf prints a log entry if quiet is false.
func (u *UnifiPoller) Logf(m string, v ...interface{}) {
	if !u.Quiet {
		_ = log.Output(2, fmt.Sprintf("[INFO] "+m, v...))
	}
}

// LogDebugf prints a debug log entry if debug is true and quite is false
func (u *UnifiPoller) LogDebugf(m string, v ...interface{}) {
	if u.Debug && !u.Quiet {
		_ = log.Output(2, fmt.Sprintf("[DEBUG] "+m, v...))
	}
}

// LogErrorf prints an error log entry.
func (u *UnifiPoller) LogErrorf(m string, v ...interface{}) {
	_ = log.Output(2, fmt.Sprintf("[ERROR] "+m, v...))
}
