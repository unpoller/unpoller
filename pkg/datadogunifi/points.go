package datadogunifi

import (
	"fmt"
	"strings"
)

func tag(name string, value any) string {
	return fmt.Sprintf("%s:%v", name, value)
}

func tagMapToTags(tagMap map[string]string) []string {
	tags := make([]string, 0)
	for k, v := range tagMap {
		tags = append(tags, tag(k, v))
	}

	return tags
}

func tagMapToSimpleStrings(tagMap map[string]string) string {
	result := ""
	for k, v := range tagMap {
		result = fmt.Sprintf("%s%s=\"%v\", ", result, k, v)
	}

	return strings.TrimRight(result, ", ")
}

func metricNamespace(namespace string) func(string) string {
	return func(name string) string {
		return fmt.Sprintf("unifi.%s.%s", namespace, name)
	}
}

func reportGaugeForFloat64Map(r report, metricName func(string) string, data map[string]float64, tags map[string]string) {
	for name, value := range data {
		_ = r.reportGauge(metricName(name), value, tagMapToTags(tags))
	}
}

// cleanTags removes any tag that is empty.
func cleanTags(tags map[string]string) map[string]string {
	for i := range tags {
		if tags[i] == "" {
			delete(tags, i)
		}
	}

	return tags
}

func boolToFloat64(v bool) float64 {
	if v {
		return 1.0
	}

	return 0.0
}
