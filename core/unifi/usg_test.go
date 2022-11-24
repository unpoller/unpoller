package unifi // nolint: testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUSGUnmarshalJSON(t *testing.T) {
	testcontroller511 := `{
  "gw": {
      "site_id": "mySite",
      "o": "gw",
      "oid": "00:00:00:00:00:00",
      "gw": "00:00:00:00:00:00",
      "time": 1577742600000,
      "datetime": "2019-12-30T09:50:00Z",
      "bytes": 0,
      "duration": 3590568000,
      "wan-rx_packets": 299729434558,
      "wan-rx_bytes": 299882768958208,
      "wan-tx_packets": 249639259523,
      "wan-tx_bytes": 169183252492369,
      "lan-rx_packets": 78912349453,
      "lan-rx_bytes": 37599596992669,
      "lan-tx_packets": 12991234992,
      "lan-tx_bytes": 11794664098210}}`

	testcontroller510 := `{
    "site_id": "mySite",
    "o": "gw",
    "oid": "00:00:00:00:00:00",
    "gw": "00:00:00:00:00:00",
    "time": 1577742600000,
    "datetime": "2019-12-30T09:50:00Z",
    "bytes": 0,
    "duration": 3590568000,
    "wan-rx_packets": 299729434558,
    "wan-rx_bytes": 299882768958208,
    "wan-tx_packets": 249639259523,
    "wan-tx_bytes": 169183252492369,
    "lan-rx_packets": 78912349453,
    "lan-rx_bytes": 37599596992669,
    "lan-tx_packets": 12991234992,
    "lan-tx_bytes": 11794664098210}`

	t.Parallel()
	a := assert.New(t)

	u := &USGStat{}
	lanRx := 37599596992669
	err := u.UnmarshalJSON([]byte(testcontroller510))
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(float64(lanRx), u.LanRxBytes.Val, "data was not properly unmarshaled")

	u = &USGStat{} // reset
	err = u.UnmarshalJSON([]byte(testcontroller511))
	a.Nil(err, "must be no error unmarshaling test strings")
	a.Equal(float64(lanRx), u.LanRxBytes.Val, "data was not properly unmarshaled")
}
