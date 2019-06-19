package unifipoller

import (
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
// It also incriments the error counter.
func (u *UnifiPoller) LogErrors(errs []error, prefix string) {
	for _, err := range errs {
		if err != nil {
			u.errorCount++
			log.Printf("[ERROR] (%v/%v) %v: %v", prefix, err.Error(), u.errorCount, u.MaxErrors)
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
		log.Printf("[INFO] "+m, v...)
	}
}
