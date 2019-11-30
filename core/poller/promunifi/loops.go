package promunifi

// This file contains all the loop methods for each device type, clients and sites.
// Moved them here to consolate clutter from the other files. Also, if these change,
// they usually all change at once since they're pretty much the same code.

func (u *unifiCollector) loopSites(r report) {
	if r.metrics() == nil || len(r.metrics().Sites) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, s := range r.metrics().Sites {
			u.exportSite(r, s)
		}
	}()
}

func (u *unifiCollector) loopUAPs(r report) {
	if r.metrics() == nil || r.metrics().Devices == nil || len(r.metrics().Devices.UAPs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range r.metrics().Devices.UAPs {
			u.exportUAP(r, d)
		}
	}()
}

func (u *unifiCollector) loopUDMs(r report) {
	if r.metrics() == nil || r.metrics().Devices == nil || len(r.metrics().Devices.UDMs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range r.metrics().Devices.UDMs {
			u.exportUDM(r, d)
		}
	}()
}

func (u *unifiCollector) loopUSGs(r report) {
	if r.metrics() == nil || r.metrics().Devices == nil || len(r.metrics().Devices.USGs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range r.metrics().Devices.USGs {
			u.exportUSG(r, d)
		}
	}()
}

func (u *unifiCollector) loopUSWs(r report) {
	if r.metrics() == nil || r.metrics().Devices == nil || len(r.metrics().Devices.USWs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range r.metrics().Devices.USWs {
			u.exportUSW(r, d)
		}
	}()
}

func (u *unifiCollector) loopClients(r report) {
	if r.metrics() == nil || len(r.metrics().Clients) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, c := range r.metrics().Clients {
			u.exportClient(r, c)
		}
	}()
}
