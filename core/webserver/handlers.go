package webserver

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

/* This file has the methods that pass out actual content. */

// Returns the main index file.
// If index.html becomes a template, this is where it can be compiled.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	index := filepath.Join(s.HTMLPath, "index.html")
	http.ServeFile(w, r, index)
}

// Arbitrary /health handler.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.handleDone(w, []byte("OK"), mimeHTML)
}

// Returns static files from static-files path. /css, /js, /img (/images, /image).
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	switch v := mux.Vars(r)["sub"]; v {
	case "image", "img":
		dir := http.Dir(filepath.Join(s.HTMLPath, "static", "images"))
		http.StripPrefix("/"+v, http.FileServer(dir)).ServeHTTP(w, r)
	default: // images, js, css, etc
		dir := http.Dir(filepath.Join(s.HTMLPath, "static", v))
		http.StripPrefix("/"+v, http.FileServer(dir)).ServeHTTP(w, r)
	}
}

// Returns web server and poller configs. /api/v1/config.
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"poller": s.Collect.Poller()}
	s.handleJSON(w, data)
}

// Returns a list of input and output plugins: /api/v1/plugins.
func (s *Server) handlePlugins(w http.ResponseWriter, r *http.Request) {
	data := map[string][]string{"inputs": s.Collect.Inputs(), "outputs": s.Collect.Outputs()}
	s.handleJSON(w, data)
}

// Returns an output plugin's data: /api/v1/output/{output}.
func (s *Server) handleOutput(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	c := s.plugins.getOutput(vars["output"])
	if c == nil {
		s.handleMissing(w, r)
		return
	}

	c.RLock()
	defer c.RUnlock()

	switch value := vars["value"]; vars["sub"] {
	default:
		s.handleJSON(w, c.Config)
	case "events":
		s.handleJSON(w, c.Events)
	case "counters":
		if value == "" {
			s.handleJSON(w, c.Counter)
		} else {
			s.handleJSON(w, map[string]int64{value: c.Counter[value]})
		}
	}
}

// Returns an input plugin's data: /api/v1/input/{input}.
func (s *Server) handleInput(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	c := s.plugins.getInput(vars["input"])
	if c == nil {
		s.handleMissing(w, r)
		return
	}

	c.RLock()
	defer c.RUnlock()

	switch value := vars["value"]; vars["sub"] {
	default:
		s.handleJSON(w, c.Config)
	case "events":
		s.handleJSON(w, c.Events)
	case "sites":
		s.handleJSON(w, c.Sites)
	case "devices":
		s.handleJSON(w, c.Devices)
	case "clients":
		s.handleJSON(w, c.Clients)
	case "counters":
		if value != "" {
			s.handleJSON(w, map[string]int64{value: c.Counter[value]})
		} else {
			s.handleJSON(w, c.Counter)
		}
	}
}
