package webserver

import (
	"strings"
	"sync"
	"time"
)

// Input is the data tracked for intput plugins.
// An input plugin should fill this data every time it polls this data.
// Partial update are OK. Set non-updated fields to nil and they're ignored.
type Input struct {
	Name         string
	Sites        Sites
	Events       Events
	Devices      Devices
	Clients      Clients
	Config       interface{}
	Counter      map[string]int64
	sync.RWMutex // Locks this data structure.
}

// Output is the data tracked for output plugins.
// Output plugins should fill this data on startup,
// and regularly update counters for things worth counting.
// Setting Config will overwrite previous value.
type Output struct {
	Name         string
	Events       Events
	Config       interface{}
	Counter      map[string]int64
	sync.RWMutex // Locks this data structure.
}

/*
These are minimal types to display a small set of data on the web interface.
These may be expanded upon, in time, as users express their needs and wants.
*/

// Sites is a list of network locations.
type Sites []*Site

// Site is a network location and its meta data.
type Site struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Source     string `json:"source"`
	Controller string `json:"controller"`
}

// Events is all the events a plugin has. string = Controller.UUID + text.
type Events map[string]*EventGroup

// EventGroup allows each plugin to have a map of events. ie. one map per controller.
type EventGroup struct {
	Latest time.Time `json:"latest"`
	Events []*Event  `json:"events"`
}

// Event is like a log message.
type Event struct {
	Ts   time.Time         `json:"ts"`
	Msg  string            `json:"msg"`
	Tags map[string]string `json:"tags,omitempty"`
}

func (e Events) Groups(filter string) (groups []string) {
	for n := range e {
		if filter == "" || strings.HasPrefix(n, filter) {
			groups = append(groups, n)
		}
	}

	return groups
}

// add adds a new event and makes sure the slice is not too big.
func (e *EventGroup) add(event *Event, max int) {
	if !e.Latest.Before(event.Ts) {
		return // Ignore older events.
	}

	e.Latest = event.Ts
	e.Events = append(e.Events, event)

	if i := len(e.Events) - max; i > 0 {
		e.Events = e.Events[i:]
	}
}

// Devices is a list of network devices and their data.
type Devices []*Device

// Device holds the data for a network device.
type Device struct {
	Clients    int         `json:"clients"`
	Uptime     int         `json:"uptime"`
	Name       string      `json:"name"`
	SiteID     string      `json:"site_id"`
	Source     string      `json:"source"`
	Controller string      `json:"controller"`
	MAC        string      `json:"mac"`
	IP         string      `json:"ip"`
	Type       string      `json:"type"`
	Model      string      `json:"model"`
	Version    string      `json:"version"`
	Config     interface{} `json:"config,omitempty"`
}

func (c Devices) Filter(siteid string) (devices []*Device) {
	for _, n := range c {
		if siteid == "" || n.SiteID == siteid {
			devices = append(devices, n)
		}
	}

	return devices
}

// Clients is a list of clients with their data.
type Clients []*Client

// Client holds the data for a network client.
type Client struct {
	Rx         int64     `json:"rx_bytes"`
	Tx         int64     `json:"tx_bytes"`
	Name       string    `json:"name"`
	SiteID     string    `json:"site_id"`
	Source     string    `json:"source"`
	Controller string    `json:"controller"`
	MAC        string    `json:"mac"`
	IP         string    `json:"ip"`
	Type       string    `json:"type"`
	DeviceMAC  string    `json:"device_mac"`
	Since      time.Time `json:"since"`
	Last       time.Time `json:"last"`
}

func (c Clients) Filter(siteid string) (clients []*Client) {
	for _, n := range c {
		if siteid == "" || n.SiteID == siteid {
			clients = append(clients, n)
		}
	}

	return clients
}
