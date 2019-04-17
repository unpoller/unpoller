# Go Library: `unifi`

It connects to a Unifi Controller, given a url, username and password. Returns
an authenticated http Client you may use to query the device for data. Also
contains some built-in methods for de-serializing common client and device
data. The data is provided in a large struct you can consume in your application.

This library also contains methods to export the Unifi data in InfluxDB format,
and this can be used as an example to base your own metrics collection methods.

If more features are requested, I'll certainly consider them. Do you need to do
more than just collect data? [Let me know](https://github.com/golift/unifi/issues/new)!
Pull requests and feedback are welcomed!

Here's a working example:
```golang
package main

import "log"
import "github.com/golift/unifi"

func main() {
	username := "admin"
	password := "superSecret1234"
	URL := "https://127.0.0.1:8443/"
	uni, err := unifi.NewUnifi(username, password, URL, false)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	// Log with log.Printf or make your own interface that accepts (msg, fmt)
	uni.ErrorLog = log.Printf
	uni.DebugLog = log.Printf

	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	clients, err := uni.GetClients(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	devices, err := uni.GetDevices(sites)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	log.Println(len(clients.UCLs), "Clients connected:")
	for i, client := range clients.UCLs {
		log.Println(i+1, client.ID, client.Hostname, client.IP, client.Name, client.LastSeen)
	}

	log.Println(len(devices.USWs), "Unifi Switches Found")
	log.Println(len(devices.USGs), "Unifi Gateways Found")

	log.Println(len(devices.UAPs), "Unifi Wireless APs Found:")
	for i, uap := range devices.UAPs {
		log.Println(i+1, uap.Name, uap.IP)
	}
}
```
