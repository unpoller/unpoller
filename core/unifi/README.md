# Go Library: `unifi`

It connects to a Unifi Controller, given a url, username and password. Returns
an authenticated http Client you may use to query the device for data. Also
contains some built-in methods for de-serializing common client and device
data. The data is provided in a large struct you can consume in your application.

This library also contains methods to export the Unifi data in InfluxDB format.

Pull requests and feedback are welcomed!


Here's a working example:
```golang
package main

import (
	"log"
	"github.com/golift/unifi"
)

func main() {
	username := "admin"
	password := "superSecret1234"
	URL := "https://127.0.0.1:8443/"
	uni, err := unifi.AuthController(username, password, URL, false)
	if err != nil {
		log.Fatalln("error:", err)
	}
  // Log with log.Printf or make your own interface.
  uni.ErrorLog = log.Printf
  uni.DebugLog = log.Printf
	clients, err := uni.GetUnifiClients()
	if err != nil {
		log.Fatalln("error:", err)
	}
	log.Println(len(clients), "clients connected:")
	for i, client := range clients {
		log.Println(i+1, client.ID, client.Hostname, client.IP, client.Name, client.LastSeen)
	}
	devices, err := uni.GetUnifiDevices()
	if err != nil {
		log.Fatalln("error:", err)
	}
	log.Println(len(devices.USWs), "Unifi Switches Found")
	log.Println(len(devices.USGs), "Unifi Gateways Found")

	log.Println(len(devices.UAPs), "Unifi Access Points Found:")
	for i, uap := range devices.UAPs {
		log.Println(i+1, uap.Name, uap.IP)
	}
}
```
