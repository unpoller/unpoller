package unifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUAPUnmarshalJSON(t *testing.T) {
	testcontroller511 := `{
	  "ap": {
	    "site_id": "mySite",
	    "o": "ap",
	    "oid": "00:00:00:00:00:00",
	    "ap": "00:00:00:00:00:00",
			"time": 1577742600000,
      "datetime": "2019-12-30T09:50:00Z",
	    "user-wifi1-rx_packets": 6596670,
	    "user-wifi0-rx_packets": 42659527,
	    "user-rx_packets": 49294197,
	    "guest-rx_packets": 0,
	    "wifi0-rx_packets": 42639527,
	    "wifi1-rx_packets": 6591670,
	    "rx_packets": 49299197}}`

	testcontroller510 := `{
		"site_id": "mySite",
		"o": "ap",
		"oid": "00:00:00:00:00:00",
		"ap": "00:00:00:00:00:00",
		"time": 1577742600000,
		"datetime": "2019-12-30T09:50:00Z",
		"user-wifi1-rx_packets": 6596670,
		"user-wifi0-rx_packets": 42659527,
		"user-rx_packets": 49294197,
		"guest-rx_packets": 0,
		"wifi0-rx_packets": 42639527,
		"wifi1-rx_packets": 6591670,
		"rx_packets": 49299197}`

	t.Parallel()
	a := assert.New(t)
	rxPakcets := 49299197
	u := &UAPStat{}
	err := u.UnmarshalJSON([]byte(testcontroller510))
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(float64(rxPakcets), u.RxPackets.Val, "data was not properly unmarshaled")

	u = &UAPStat{} // reset
	err = u.UnmarshalJSON([]byte(testcontroller511))
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(float64(rxPakcets), u.RxPackets.Val, "data was not properly unmarshaled")
}
