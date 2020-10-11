package datadogunifi

import (
	"fmt"
)

func tag(name string, value interface{}) string {
	return fmt.Sprintf("%s:%v", name, value)
}

func metricNamespace(namespace string) func(string) string {
	return func(name string) string {
		return fmt.Sprintf("%s.%s", namespace, name)
	}
}

func reportGaugeForMap(r report, metricName func(string) string, data map[string]float64, tags []string) {
	for name, value := range data {
		r.reportGauge(metricName(name), value, tags)
	}
}
