package webserver

import (
	"net/http"
	"path/filepath"
	"time"

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

// Returns poller configs and/or plugins. /api/v1/config.
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch vars["sub"] {
	case "":
		data := map[string]interface{}{
			"inputs":  s.Collect.Inputs(),
			"outputs": s.Collect.Outputs(),
			"poller":  s.Collect.Poller(),
			"uptime":  int(time.Since(s.start).Round(time.Second).Seconds()),
		}
		s.handleJSON(w, data)
	case "plugins":
		data := map[string][]string{
			"inputs":  s.Collect.Inputs(),
			"outputs": s.Collect.Outputs(),
		}
		s.handleJSON(w, data)
	default:
		s.handleMissing(w, r)
	}
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

	switch val := vars["value"]; vars["sub"] {
	default:
		s.handleJSON(w, c.Config)
	case "eventgroups":
		s.handleJSON(w, c.Events.Groups(val))
	case "events":
		switch events, ok := c.Events[val]; {
		case val == "":
			s.handleJSON(w, c.Events)
		case ok:
			s.handleJSON(w, events)
		default:
			s.handleMissing(w, r)
		}
	case "counters":
		if val == "" {
			s.handleJSON(w, c.Counter)
		} else {
			s.handleJSON(w, map[string]int64{val: c.Counter[val]})
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

	switch val := vars["value"]; vars["sub"] {
	default:
		s.handleJSON(w, c.Config)
	case "eventgroups":
		s.handleJSON(w, c.Events.Groups(val))
	case "events":
		switch events, ok := c.Events[val]; {
		case val == "":
			s.handleJSON(w, c.Events)
		case ok:
			s.handleJSON(w, events)
		default:
			s.handleMissing(w, r)
		}
	case "sites":
		s.handleJSON(w, c.Sites)
	case "devices":
		s.handleJSON(w, c.Devices.Filter(val))
	case "clients":
		s.handleJSON(w, c.Clients.Filter(val))
	case "counters":
		if val != "" {
			s.handleJSON(w, map[string]int64{val: c.Counter[val]})
		} else {
			s.handleJSON(w, c.Counter)
		}
	}
}
