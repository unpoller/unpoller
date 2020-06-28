module github.com/unifi-poller/inputunifi

go 1.14

replace github.com/unifi-poller/webserver => ../webserver

replace github.com/unifi-poller/poller => ../poller

require (
	github.com/pkg/errors v0.9.1
	github.com/unifi-poller/poller v0.0.8-0.20200626082958-a9a7092a5684
	github.com/unifi-poller/unifi v0.0.6-0.20200625090439-421046871a37
	github.com/unifi-poller/webserver v0.0.0-20200628114213-2b89a50ff1c0
)
