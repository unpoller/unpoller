package main

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

// logErrors writes a slice of errors, with a prefix, to log-out.
func logErrors(errs []error, prefix string) {
	for _, err := range errs {
		if err != nil {
			log.Println("[ERROR]", prefix+":", err.Error())
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
func (c *Config) Logf(m string, v ...interface{}) {
	if !c.Quiet {
		log.Printf("[INFO] "+m, v...)
	}
}
