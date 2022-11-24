package unifi

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// anomaly is the type UniFi returns, but not the type this library returns.
type anomaly struct {
	Anomaly    string  `json:"anomaly"`
	MAC        string  `json:"mac"`
	Timestamps []int64 `json:"timestamps"`
}

// Anomaly is the reformatted data type that this library returns.
type Anomaly struct {
	Datetime   time.Time
	SourceName string
	SiteName   string
	Anomaly    string
	DeviceMAC  string
	//	DeviceName string // we do not have this....
}

// GetAnomalies returns Anomalies for a list of Sites.
func (u *Unifi) GetAnomalies(sites []*Site, timeRange ...time.Time) ([]*Anomaly, error) {
	data := []*Anomaly{}

	for _, site := range sites {
		response, err := u.GetAnomaliesSite(site, timeRange...)
		if err != nil {
			return data, err
		}

		data = append(data, response...)
	}

	return data, nil
}

// GetAnomaliesSite retreives the Anomalies for a single Site.
func (u *Unifi) GetAnomaliesSite(site *Site, timeRange ...time.Time) ([]*Anomaly, error) {
	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for Anomalies, site %s", site.SiteName)

	var (
		path      = fmt.Sprintf(APIAnomaliesPath, site.Name)
		anomalies = anomalies{}
		data      struct {
			Data []*anomaly `json:"data"`
		}
	)

	if params, err := makeAnomalyParams("hourly", timeRange...); err != nil {
		return anomalies, err
	} else if err := u.GetData(path+params, &data, ""); err != nil {
		return anomalies, err
	}

	for _, d := range data.Data {
		for _, ts := range d.Timestamps {
			anomalies = append(anomalies, &Anomaly{
				Datetime:   time.Unix(ts/int64(time.Microsecond), 0),
				SourceName: u.URL,
				SiteName:   site.SiteName,
				Anomaly:    d.Anomaly,
				DeviceMAC:  d.MAC,
				//				DeviceName: d.Anomaly,
			})
		}
	}

	sort.Sort(anomalies)

	return anomalies, nil
}

// anomalies satisfies the sort.Sort interface.
type anomalies []*Anomaly

// Len satisfies sort.Interface.
func (a anomalies) Len() int {
	return len(a)
}

// Swap satisfies sort.Interface.
func (a anomalies) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less satisfies sort.Interface. Sort our list by Datetime.
func (a anomalies) Less(i, j int) bool {
	return a[i].Datetime.Before(a[j].Datetime)
}

func makeAnomalyParams(scale string, timeRange ...time.Time) (string, error) {
	out := []string{}

	if scale != "" {
		out = append(out, "scale="+scale)
	}

	switch len(timeRange) {
	case 0:
		end := time.Now().Unix() * int64(time.Microsecond)
		out = append(out, "end="+strconv.FormatInt(end, 10))
	case 1:
		start := timeRange[0].Unix() * int64(time.Microsecond)
		end := time.Now().Unix() * int64(time.Microsecond)
		out = append(out, "end="+strconv.FormatInt(end, 10), "start="+strconv.FormatInt(start, 10))
	case 2: // nolint: gomnd
		start := timeRange[0].Unix() * int64(time.Microsecond)
		end := timeRange[1].Unix() * int64(time.Microsecond)
		out = append(out, "end="+strconv.FormatInt(end, 10), "start="+strconv.FormatInt(start, 10))
	default:
		return "", ErrInvalidTimeRange
	}

	if len(out) == 0 {
		return "", nil
	}

	return "?" + strings.Join(out, "&"), nil
}
