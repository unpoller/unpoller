package webserver

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/unifi-poller/poller"
)

/* This file has the methods that help the content-methods. Shared helpers. */

const (
	xPollerError = "X-Poller-Error"
	mimeJSON     = "application/json"
	mimeHTML     = "text/plain; charset=utf-8"
)

// basicAuth wraps web requests with simple auth (and logging).
// Called on nearly every request.
func (s *Server) basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return s.handleLog(func(w http.ResponseWriter, r *http.Request) {
		if s.Accounts.PasswordIsCorrect(r.BasicAuth()) {
			handler(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Enter Name and Password to Login!"`)
		w.WriteHeader(http.StatusUnauthorized)
	})
}

// handleLog writes an Apache-like log line. Called on every request.
func (s *Server) handleLog(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Scheme = "https"; r.TLS == nil {
			r.URL.Scheme = "http" // Set schema early in case another handler uses it.
		}

		// Use custom ResponseWriter to catch and log response data.
		response := &ResponseWriter{Writer: w, Start: time.Now()}
		handler(response, r) // Run provided handler with custom ResponseWriter.

		user, _, _ := r.BasicAuth()
		if user == "" {
			user = "-" // Only used for logs.
		}

		logf := s.Logf // Standard log.
		if response.Error != "" {
			logf = s.LogErrorf // Format an error log.
			response.Error = ` "` + response.Error + `"`
		}

		remote, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			remote = r.RemoteAddr
		}

		logf(`%s %s %s [%v] "%s %s://%s%s %s" %d %d "%s" "%s" %v%s`, remote, poller.AppName,
			user, response.Start.Format("01/02/2006:15:04:05 -07:00"), r.Method, r.URL.Scheme,
			r.Host, r.RequestURI, r.Proto, response.Code, response.Size, r.Referer(),
			r.UserAgent(), time.Since(response.Start).Round(time.Microsecond), response.Error)
	}
}

// handleMissing returns a blank 404.
func (s *Server) handleMissing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", mimeHTML)
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("404 page not found\n"))
}

// handleError is a pass-off function when a request returns an error.
func (s *Server) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", mimeHTML)
	w.Header().Set(xPollerError, err.Error()) // signal
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error() + "\n"))
}

// handleDone is a pass-off function to finish a request.
func (s *Server) handleDone(w http.ResponseWriter, b []byte, cType string) {
	w.Header().Set("Content-Type", cType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(b, []byte("\n")...))
}

// handleJSON sends a json-formatted data reply.
func (s *Server) handleJSON(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		s.handleError(w, err)
		return
	}

	s.handleDone(w, b, mimeJSON)
}

/* Custom http.ResponseWriter interface method and struct overrides. */

// ResponseWriter is used to override http.ResponseWriter in our http.FileServer.
// This allows us to catch and log the response code, size and error; maybe others.
type ResponseWriter struct {
	Code   int
	Size   int
	Error  string
	Start  time.Time
	Writer http.ResponseWriter
}

// Header sends a header to a client. Satisfies http.ResponseWriter interface.
func (w *ResponseWriter) Header() http.Header {
	return w.Writer.Header()
}

// Write sends bytes to the client. Satisfies http.ResponseWriter interface.
// This also adds the written byte count to our size total.
func (w *ResponseWriter) Write(b []byte) (int, error) {
	size, err := w.Writer.Write(b)
	w.Size += size

	return size, err
}

// WriteHeader sends an http StatusCode to a client. Satisfies http.ResponseWriter interface.
// This custom override method also saves the status code, and any error message (for logs).
func (w *ResponseWriter) WriteHeader(code int) {
	w.Error = w.Header().Get(xPollerError) // Catch and save any response error.
	w.Header().Del(xPollerError)           // Delete the temporary signal header.
	w.Code = code                          // Save the status code.
	w.Writer.WriteHeader(code)             // Pass the request through.
}
