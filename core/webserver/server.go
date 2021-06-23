// Package webserver is a UniFi Poller plugin that exports running data to a web interface.
package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/unpoller/poller"
	"golang.org/x/crypto/bcrypt"
)

const (
	// PluginName identifies this output plugin.
	PluginName = "WebServer"
	// DefaultPort is the default web http port.
	DefaultPort = 37288
	// DefaultEvents is the default number of events stored per plugin.
	DefaultEvents = 200
)

// Config is the webserver library input config.
type Config struct {
	Enable     bool     `json:"enable" toml:"enable" xml:"enable,attr" yaml:"enable"`
	SSLCrtPath string   `json:"ssl_cert_path" toml:"ssl_cert_path" xml:"ssl_cert_path" yaml:"ssl_cert_path"`
	SSLKeyPath string   `json:"ssl_key_path" toml:"ssl_key_path" xml:"ssl_key_path" yaml:"ssl_key_path"`
	Port       uint     `json:"port" toml:"port" xml:"port" yaml:"port"`
	Accounts   accounts `json:"accounts" toml:"accounts" xml:"accounts" yaml:"accounts"`
	HTMLPath   string   `json:"html_path" toml:"html_path" xml:"html_path" yaml:"html_path"`
	MaxEvents  uint     `json:"max_events" toml:"max_events" xml:"max_events" yaml:"max_events"`
}

// accounts stores a map of usernames and password hashes.
type accounts map[string]string

// Server is the main library struct/data.
type Server struct {
	*Config `json:"webserver" toml:"webserver" xml:"webserver" yaml:"webserver"`
	server  *http.Server
	plugins *webPlugins
	Collect poller.Collect
	start   time.Time
}

// init is how this modular code is initialized by the main app.
// This module adds itself as an output module to the poller core.
func init() { // nolint: gochecknoinits
	s := &Server{plugins: plugins, start: time.Now(), Config: &Config{
		Port:      DefaultPort,
		HTMLPath:  filepath.Join(poller.DefaultObjPath, "web"),
		MaxEvents: DefaultEvents,
	}}
	plugins.Config = s.Config

	poller.NewOutput(&poller.Output{
		Name:   PluginName,
		Config: s,
		Method: s.Run,
	})
}

// Run starts the server and gets things going.
func (s *Server) Run(c poller.Collect) error {
	if s.Collect = c; s.Config == nil || s.Port == 0 || s.HTMLPath == "" || !s.Enable {
		s.Logf("Internal web server disabled!")
		return nil
	}

	if _, err := os.Stat(s.HTMLPath); err != nil {
		return fmt.Errorf("problem with HTML path: %w", err)
	}

	UpdateOutput(&Output{Name: PluginName, Config: s.Config})

	return s.Start()
}

// Start gets the web server going.
func (s *Server) Start() (err error) {
	s.server = &http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(int(s.Port)),
		WriteTimeout: time.Minute,
		ReadTimeout:  time.Minute,
		IdleTimeout:  time.Minute,
		Handler:      s.newRouter(), // *mux.Router
	}

	if s.SSLCrtPath == "" || s.SSLKeyPath == "" {
		s.Logf("Web server starting without SSL. Listening on HTTP port %d", s.Port)
		err = s.server.ListenAndServe()
	} else {
		s.Logf("Web server starting with SSL. Listening on HTTPS port %d", s.Port)
		err = s.server.ListenAndServeTLS(s.SSLCrtPath, s.SSLKeyPath)
	}

	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("web server: %w", err)
	}

	return nil
}

func (s *Server) newRouter() *mux.Router {
	router := mux.NewRouter()
	// special routes
	router.Handle("/debug/vars", http.DefaultServeMux).Methods("GET")        // unauthenticated expvar
	router.HandleFunc("/health", s.handleLog(s.handleHealth)).Methods("GET") // unauthenticated health
	// main web app/files/js/css
	router.HandleFunc("/", s.basicAuth(s.handleIndex)).Methods("GET", "POST")
	router.PathPrefix("/{sub:css|js|img|image|images}/").Handler((s.basicAuth(s.handleStatic))).Methods("GET")
	// api paths for json dumps
	router.HandleFunc("/api/v1/config", s.basicAuth(s.handleConfig)).Methods("GET")
	router.HandleFunc("/api/v1/config/{sub}", s.basicAuth(s.handleConfig)).Methods("GET")
	router.HandleFunc("/api/v1/config/{sub}/{value}", s.basicAuth(s.handleConfig)).Methods("GET", "POST")
	router.HandleFunc("/api/v1/input/{input}", s.basicAuth(s.handleInput)).Methods("GET")
	router.HandleFunc("/api/v1/input/{input}/{sub}", s.basicAuth(s.handleInput)).Methods("GET")
	router.HandleFunc("/api/v1/input/{input}/{sub}/{value}", s.basicAuth(s.handleInput)).Methods("GET", "POST")
	router.HandleFunc("/api/v1/output/{output}", s.basicAuth(s.handleOutput)).Methods("GET")
	router.HandleFunc("/api/v1/output/{output}/{sub}", s.basicAuth(s.handleOutput)).Methods("GET")
	router.HandleFunc("/api/v1/output/{output}/{sub}/{value}", s.basicAuth(s.handleOutput)).Methods("GET", "POST")
	router.PathPrefix("/").Handler(s.basicAuth(s.handleMissing)).Methods("GET", "POST", "PUT") // 404 everything.

	return router
}

// PasswordIsCorrect returns true if the provided password matches a user's account.
func (a accounts) PasswordIsCorrect(user, pass string, ok bool) bool {
	if len(a) == 0 {
		return true // No accounts defined in config; allow anyone.
	} else if !ok {
		return false // r.BasicAuth() failed, not a valid user.
	} else if user, ok = a[user]; !ok { // The user var is now the password hash.
		return false // The username provided doesn't exist.
	}

	// If this is returns nil, the provided password matches, so return true.
	return bcrypt.CompareHashAndPassword([]byte(user), []byte(pass)) == nil
}
