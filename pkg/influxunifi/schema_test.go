package influxunifi_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestInfluxSchemaNoTagFieldOverlap(t *testing.T) {
	t.Helper()

	yamlFile, err := os.ReadFile("integration_test_expectations.yaml")
	require.NoError(t, err)

	var data testExpectations

	err = yaml.Unmarshal(yamlFile, &data)
	require.NoError(t, err)

	for measurement, spec := range data.Points {
		tags := make(map[string]struct{}, len(spec.Tags))
		for _, tag := range spec.Tags {
			tags[tag] = struct{}{}
		}

		for field := range spec.Fields {
			if _, ok := tags[field]; ok {
				t.Errorf("measurement %q has overlapping tag/field key %q", measurement, field)
			}
		}
	}
}
