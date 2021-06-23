# Go Library: `unifi`

It connects to a Unifi Controller, given a url, username and password. Returns
an authenticated http Client you may use to query the device for data. Also
contains some built-in methods for de-serializing common client and device
data. The data is provided in a large struct you can consume in your application.

This library is designed to PULL data FROM the controller. It has no methods that
update settings or change things on the controller.
[Someone expressed interest](https://github.com/unpoller/unifi/issues/31) in
adding methods to update data, and I'm okay with that. I'll even help add them.
[Tell me what you want to do](https://github.com/unpoller/unifi/issues/new), and we'll make it happen.

Pull requests, feature requests, code reviews and feedback are welcomed!

Here's a working example:
```golang
package main

import "log"
import "github.com/unpoller/unifi"

func main() {
	c := *unifi.Config{
		User: "admin",
		Pass: "superSecret1234",
		URL:  "https://127.0.0.1:8443/",
		// Log with log.Printf or make your own interface that accepts (msg, fmt)
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}
	uni, err := unifi.NewUnifi(c)
	if err != nil {
		log.Fatalln("Error:", err)
	}

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

	log.Println(len(sites), "Unifi Sites Found: ", sites)
	log.Println(len(clients), "Clients connected:")
	for i, client := range clients {
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
