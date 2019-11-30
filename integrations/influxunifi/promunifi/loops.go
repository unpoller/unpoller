package promunifi

// This file contains all the loop methods for each device type, clients and sites.
// Moved them here to consolate clutter from the other files. Also, if these change,
// they usually all change at once since they're pretty much the same code.

func (u *unifiCollector) loopSites(r report) {
	m := r.metrics()
	if m == nil || len(m.Sites) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, s := range m.Sites {
			u.exportSite(r, s)
		}
	}()
}

func (u *unifiCollector) loopUAPs(r report) {
	m := r.metrics()
	if m == nil || m.Devices == nil || len(m.Devices.UAPs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range m.Devices.UAPs {
			u.exportUAP(r, d)
		}
	}()
}

func (u *unifiCollector) loopUDMs(r report) {
	m := r.metrics()
	if m == nil || m.Devices == nil || len(m.Devices.UDMs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range m.Devices.UDMs {
			u.exportUDM(r, d)
		}
	}()
}

func (u *unifiCollector) loopUSGs(r report) {
	m := r.metrics()
	if m == nil || m.Devices == nil || len(m.Devices.USGs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range m.Devices.USGs {
			u.exportUSG(r, d)
		}
	}()
}

func (u *unifiCollector) loopUSWs(r report) {
	m := r.metrics()
	if m == nil || m.Devices == nil || len(m.Devices.USWs) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, d := range m.Devices.USWs {
			u.exportUSW(r, d)
		}
	}()
}

func (u *unifiCollector) loopClients(r report) {
	m := r.metrics()
	if m == nil || len(m.Clients) < 1 {
		return
	}
	r.add()
	go func() {
		defer r.done()
		for _, c := range m.Clients {
			u.exportClient(r, c)
		}
	}()
}
