package promunifi

// This file contains all the loop methods for each device type, clients and sites.
// Moved them here to consolate clutter from the other files. Also, if these change,
// they usually all change at once since they're pretty much the same code.

func (u *promUnifi) loopSites(r report) {
	defer r.done()
	for _, s := range r.metrics().Sites {
		u.exportSite(r, s)
	}
}

func (u *promUnifi) loopUAPs(r report) {
	defer r.done()
	for _, d := range r.metrics().UAPs {
		u.exportUAP(r, d)
	}
}

func (u *promUnifi) loopUDMs(r report) {
	defer r.done()
	for _, d := range r.metrics().UDMs {
		u.exportUDM(r, d)
	}
}

func (u *promUnifi) loopUSGs(r report) {
	defer r.done()
	for _, d := range r.metrics().USGs {
		u.exportUSG(r, d)
	}
}

func (u *promUnifi) loopUSWs(r report) {
	defer r.done()
	for _, d := range r.metrics().USWs {
		u.exportUSW(r, d)
	}
}

func (u *promUnifi) loopClients(r report) {
	defer r.done()
	for _, c := range r.metrics().Clients {
		u.exportClient(r, c)
	}
}
