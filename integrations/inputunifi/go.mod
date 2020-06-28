module github.com/unifi-poller/inputunifi

go 1.14

replace github.com/unifi-poller/webserver => ../webserver

replace github.com/unifi-poller/poller => ../poller

require (
	github.com/pkg/errors v0.9.1
	github.com/unifi-poller/poller v0.0.8-0.20200628131550-26430cac16c1
	github.com/unifi-poller/unifi v0.0.6-0.20200628131946-cd307d4b367b
	github.com/unifi-poller/webserver v0.0.0-20200628132023-9a5dfcd56166
)
