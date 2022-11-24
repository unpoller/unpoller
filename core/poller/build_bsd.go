// +build darwin freebsd netbsd openbsd

package poller

// DefaultConfFile is where to find config if --config is not prvided.
const DefaultConfFile = "/etc/unifi-poller/up.conf,/usr/local/etc/unifi-poller/up.conf"

// DefaultObjPath is the path to look for shared object libraries (plugins).
const DefaultObjPath = "/usr/local/lib/unifi-poller"
