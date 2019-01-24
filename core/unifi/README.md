# Go Library: `unifi`

It connects to a Unifi Controller, given a url, username and password. Returns
an authenticated http Client you may use to query the device for data. Also
contains some built-in methods for de-serializing common client and device
data. The included asset interface currently only works for InfluxDB but could
probably be modified to support other output mechanisms; not sure really.

Pull requests and feedback are welcomed!

This lib is rudimentary and gets a job done for the tool at hand. It could be
used to base your own library. Good luck!
