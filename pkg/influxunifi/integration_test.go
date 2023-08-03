package influxunifi_test

import (
	"fmt"
	"log"
	"os"
	"sort"
	"testing"
	"time"

	influxV1Models "github.com/influxdata/influxdb1-client/models"
	influxV1 "github.com/influxdata/influxdb1-client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unpoller/unpoller/pkg/influxunifi"
	"github.com/unpoller/unpoller/pkg/unittest"
	"golift.io/cnfg"
	"gopkg.in/yaml.v3"
)

var errNotImplemented = fmt.Errorf("not implemented")

type mockInfluxV1Client struct {
	databases map[string]bool
	points    map[string]foundPoint
}

type foundPoint struct {
	Tags   map[string]string `json:"tags"`
	Fields map[string]string `json:"fields"`
}

func newMockInfluxV1Client() *mockInfluxV1Client {
	return &mockInfluxV1Client{
		databases: make(map[string]bool, 0),
		points:    make(map[string]foundPoint, 0),
	}
}

func (m *mockInfluxV1Client) toTestData() testExpectations {
	dbs := make([]string, 0)
	for k := range m.databases {
		dbs = append(dbs, k)
	}

	sort.Strings(dbs)

	result := testExpectations{
		Databases: dbs,
		Points:    map[string]testPointExpectation{},
	}

	for k, p := range m.points {
		tags := make([]string, 0)
		for t := range p.Tags {
			tags = append(tags, t)
		}

		sort.Strings(tags)

		result.Points[k] = testPointExpectation{
			Tags:   tags,
			Fields: p.Fields,
		}
	}

	return result
}

func (m *mockInfluxV1Client) Ping(_ time.Duration) (time.Duration, string, error) {
	return time.Millisecond, "", nil
}

func influxV1FieldTypeToString(ft influxV1Models.FieldType) string {
	switch ft {
	case influxV1Models.Integer:
		return "int"
	case influxV1Models.Float:
		return "float"
	case influxV1Models.Boolean:
		return "bool"
	case influxV1Models.Empty:
		return "none"
	case influxV1Models.String:
		return "string"
	case influxV1Models.Unsigned:
		return "uint"
	default:
		return "unknown"
	}
}

func (m *mockInfluxV1Client) Write(bp influxV1.BatchPoints) error {
	m.databases[bp.Database()] = true
	for _, p := range bp.Points() {
		if existing, ok := m.points[p.Name()]; !ok {
			m.points[p.Name()] = foundPoint{
				Tags:   p.Tags(),
				Fields: make(map[string]string),
			}
		} else {
			for k, v := range p.Tags() {
				existing.Tags[k] = v
			}
		}

		fields, err := p.Fields()
		if err != nil {
			continue
		}

		point, _ := influxV1Models.NewPoint(p.Name(), influxV1Models.NewTags(p.Tags()), fields, p.Time())

		fieldIter := point.FieldIterator()
		for fieldIter.Next() {
			fieldName := string(fieldIter.FieldKey())
			fieldType := influxV1FieldTypeToString(fieldIter.Type())

			if _, exists := m.points[p.Name()].Fields[fieldName]; exists {
				if fieldType != "" {
					m.points[p.Name()].Fields[fieldName] = fieldType
				}
			} else {
				if fieldType == "" {
					m.points[p.Name()].Fields[fieldName] = "unknown"
				} else {
					m.points[p.Name()].Fields[fieldName] = fieldType
				}
			}
		}
	}

	return nil
}

func (m *mockInfluxV1Client) Query(_ influxV1.Query) (*influxV1.Response, error) {
	return nil, errNotImplemented
}

func (m *mockInfluxV1Client) QueryAsChunk(_ influxV1.Query) (*influxV1.ChunkedResponse, error) {
	return nil, errNotImplemented
}

func (m *mockInfluxV1Client) Close() error {
	return nil
}

type testPointExpectation struct {
	Tags   []string          `json:"tags"`
	Fields map[string]string `json:"fields"`
}

type testExpectations struct {
	Databases []string                        `json:"databases"`
	Points    map[string]testPointExpectation `json:"points"`
}

func TestInfluxV1Integration(t *testing.T) {
	// load test expectations file
	yamlFile, err := os.ReadFile("integration_test_expectations.yaml")
	require.NoError(t, err)

	var testExpectationsData testExpectations
	err = yaml.Unmarshal(yamlFile, &testExpectationsData)
	require.NoError(t, err)

	testRig := unittest.NewTestSetup(t)
	defer testRig.Close()

	mockCapture := newMockInfluxV1Client()

	u := influxunifi.InfluxUnifi{
		Collector:      testRig.Collector,
		IsVersion2:     false,
		InfluxV1Client: mockCapture,
		InfluxDB: &influxunifi.InfluxDB{
			Config: &influxunifi.Config{
				DB:       "unpoller",
				URL:      testRig.MockServer.Server.URL,
				Interval: cnfg.Duration{Duration: time.Hour},
			},
		},
	}

	testRig.Initialize()

	u.Poll(time.Minute)

	// databases
	assert.Len(t, mockCapture.databases, 1)

	expectedKeys := unittest.NewSetFromSlice[string](testExpectationsData.Databases)
	foundKeys := unittest.NewSetFromMap[string](mockCapture.databases)
	additions, deletions := expectedKeys.Difference(foundKeys)
	assert.Len(t, additions, 0)
	assert.Len(t, deletions, 0)

	// point names
	assert.Len(t, unittest.NewSetFromMap[string](mockCapture.points).Slice(), len(testExpectationsData.Points))
	expectedKeys = unittest.NewSetFromMap[string](testExpectationsData.Points)
	foundKeys = unittest.NewSetFromMap[string](mockCapture.points)
	additions, deletions = expectedKeys.Difference(foundKeys)
	assert.Len(t, additions, 0)
	assert.Len(t, deletions, 0)

	// validate tags and fields per point
	pointNames := unittest.NewSetFromMap[string](testExpectationsData.Points).Slice()
	sort.Strings(pointNames)

	for _, pointName := range pointNames {
		expectedContent := testExpectationsData.Points[pointName]
		foundContent := mockCapture.points[pointName]
		// check tags left intact
		expectedKeys = unittest.NewSetFromSlice[string](expectedContent.Tags)
		foundKeys = unittest.NewSetFromMap[string](foundContent.Tags)
		additions, deletions = expectedKeys.Difference(foundKeys)
		assert.Len(t, additions, 0, "point \"%s\" found the following tag keys have a difference!: additions=%+v", pointName, additions)
		assert.Len(t, deletions, 0, "point \"%s\" found the following tag keys have a difference!: deletions=%+v", pointName, deletions)

		// check field keys intact
		expectedKeys = unittest.NewSetFromMap[string](expectedContent.Fields)
		foundKeys = unittest.NewSetFromMap[string](foundContent.Fields)
		additions, deletions = expectedKeys.Difference(foundKeys)
		assert.Len(
			t,
			additions,
			0,
			"point \"%s\" found the following field keys have a difference!: additions=%+v foundKeys=%+v",
			pointName,
			additions,
			foundKeys.Slice(),
		)
		assert.Len(
			t,
			deletions,
			0,
			"point \"%s\" found the following field keys have a difference!: deletions=%+v foundKeys=%+v",
			pointName,
			deletions,
			foundKeys.Slice(),
		)

		// check field types
		fieldNames := unittest.NewSetFromMap[string](expectedContent.Fields).Slice()
		sort.Strings(fieldNames)

		for _, fieldName := range fieldNames {
			expectedFieldType := expectedContent.Fields[fieldName]
			foundFieldType := foundContent.Fields[fieldName]
			assert.Equal(t, expectedFieldType, foundFieldType, "point \"%s\" field \"%s\" had a difference in declared type \"%s\" vs \"%s\", this is not safe for backwards compatibility", pointName, fieldName, expectedFieldType, foundFieldType)
		}
	}

	capturedTestData := mockCapture.toTestData()
	buf, _ := yaml.Marshal(&capturedTestData)
	log.Println("generated expectation yaml:\n" + string(buf))
}
