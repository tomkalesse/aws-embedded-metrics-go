package context

import (
	"log"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/utils"
)

type MetricsContext struct {
	Namespace                  string
	Properties                 map[string]string
	Metrics                    map[string]MetricsValue
	Meta                       map[string]string
	Dimensions                 []map[string]string
	DefaultDimensions          map[string]string
	ShouldUseDefaultDimensions bool
	Timestamp                  int64
	MetricNameAndResolutionMap map[string]utils.StorageResolution
}

type MetricsValue struct {
	Values            []float64
	Unit              utils.Unit
	StorageResolution utils.StorageResolution
}

func (m *MetricsValue) addValue(value float64) {
	m.Values = append(m.Values, value)
}

func Empty() MetricsContext {
	return MetricsContext{
		Namespace:                  "",
		Properties:                 make(map[string]string),
		Metrics:                    make(map[string]MetricsValue),
		Meta:                       make(map[string]string),
		Dimensions:                 make([]map[string]string, 0),
		Timestamp:                  0,
		MetricNameAndResolutionMap: make(map[string]utils.StorageResolution),
	}
}

func (m *MetricsContext) SetNamespace(namespace string) error {
	err := validateNamespace(namespace)
	if err != nil {
		return err
	} else {
		m.Namespace = namespace
		return nil
	}

}

func (m *MetricsContext) SetProperty(key string, value string) {
	m.Properties[key] = value
}

func (m *MetricsContext) SetTimestamp(timestamp int64) {
	m.Timestamp = timestamp
}

func (m *MetricsContext) SetDefaultDimensions(dimensions map[string]string) {
	log.Println("Received default dimensions")
	m.DefaultDimensions = dimensions
}

func (m *MetricsContext) PutDimensions(incomingDimensionSet map[string]string) error {
	err := validateDimensionSet(incomingDimensionSet)
	if err != nil {
		return err
	}

	incomingDimensionSetKeys := getMapKeys(incomingDimensionSet)

	var filteredDimensions []map[string]string
	for _, existingDimensionSet := range m.Dimensions {
		existingDimensionSetKeys := getMapKeys(existingDimensionSet)

		if len(existingDimensionSetKeys) != len(incomingDimensionSetKeys) ||
			!areSlicesEqual(existingDimensionSetKeys, incomingDimensionSetKeys) {
			filteredDimensions = append(filteredDimensions, existingDimensionSet)
		}
	}

	m.Dimensions = append(filteredDimensions, incomingDimensionSet)

	return nil

}

func (m *MetricsContext) SetDimensions(dimensionSets []map[string]string, useDefault bool) error {
	for _, dimensionSet := range dimensionSets {
		err := validateDimensionSet(dimensionSet)
		if err != nil {
			return err
		}
	}
	m.ShouldUseDefaultDimensions = useDefault
	m.Dimensions = dimensionSets
	return nil
}

func (m *MetricsContext) ResetDimensions(useDefault bool) {
	m.ShouldUseDefaultDimensions = useDefault
	m.Dimensions = make([]map[string]string, 0)
}

func (m *MetricsContext) GetDimensions() []map[string]string {
	if !m.ShouldUseDefaultDimensions {
		return m.Dimensions
	}

	if len(m.DefaultDimensions) == 0 {
		return m.Dimensions
	}

	if len(m.Dimensions) == 0 {
		return []map[string]string{m.DefaultDimensions}
	}

	mergedDimensions := make([]map[string]string, len(m.Dimensions))
	for _, customDimension := range m.Dimensions {
		mergedDimensions = append(mergedDimensions, mergeMaps(m.DefaultDimensions, customDimension))
	}

	return mergedDimensions
}

func (m *MetricsContext) PutMetric(key string, value float64, unit utils.Unit, storageResolution utils.StorageResolution) error {
	err := validateMetric(key, unit, storageResolution, m.MetricNameAndResolutionMap)
	if err != nil {
		return err
	}
	currentMetric := m.Metrics[key]
	if currentMetric.Values != nil && currentMetric.Unit != "" && currentMetric.StorageResolution != 0 {
		currentMetric.addValue(value)
	} else {
		m.Metrics[key] = MetricsValue{
			Values:            []float64{value},
			Unit:              unit,
			StorageResolution: storageResolution,
		}
	}
	m.MetricNameAndResolutionMap[key] = storageResolution
	return nil
}

func (m *MetricsContext) CreateCopyWithContext(preserveDimensions bool) MetricsContext {
	return MetricsContext{
		Namespace:                  m.Namespace,
		Properties:                 m.Properties,
		Dimensions:                 m.Dimensions,
		DefaultDimensions:          m.DefaultDimensions,
		ShouldUseDefaultDimensions: preserveDimensions,
		Timestamp:                  m.Timestamp,
	}
}

func getMapKeys(dimensionSet map[string]string) []string {
	keys := make([]string, 0, len(dimensionSet))
	for key := range dimensionSet {
		keys = append(keys, key)
	}
	return keys
}

func areSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]bool)
	for _, item := range a {
		aMap[item] = true
	}
	for _, item := range b {
		if !aMap[item] {
			return false
		}
	}
	return true
}

func mergeMaps(defaults, custom map[string]string) map[string]string {
	merged := make(map[string]string)

	for key, value := range defaults {
		merged[key] = value
	}

	for key, value := range custom {
		merged[key] = value
	}

	return merged
}
