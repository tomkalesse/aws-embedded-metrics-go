package context

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/utils"
)

type MetricsContext struct {
	Namespace                  string
	Properties                 map[string]string
	Metrics                    map[string]MetricsValue
	Meta                       map[string]any
	dimensions                 []map[string]string
	defaultDimensions          map[string]string
	shouldUseDefaultDimensions bool
	timestamp                  int64
	metricNameAndResolutionMap map[string]utils.StorageResolution
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
		Namespace:                  config.GetConfig().Namespace,
		Properties:                 make(map[string]string),
		Metrics:                    make(map[string]MetricsValue),
		Meta:                       map[string]any{"Timestamp": resolveMetaTimestamp(0)},
		dimensions:                 make([]map[string]string, 0),
		shouldUseDefaultDimensions: true,
		timestamp:                  0,
		metricNameAndResolutionMap: make(map[string]utils.StorageResolution),
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

func (m *MetricsContext) SetTimestamp(timestamp int64) error {
	err := validateTimestamp(timestamp)
	if err != nil {
		return err
	}
	m.timestamp = timestamp
	m.Meta["Timestamp"] = resolveMetaTimestamp(timestamp)
	return nil
}

func (m *MetricsContext) SetDefaultDimensions(dimensions map[string]string) {
	m.defaultDimensions = dimensions
}

func (m *MetricsContext) PutDimensions(incomingDimensionSet map[string]string) error {
	err := validateDimensionSet(incomingDimensionSet)
	if err != nil {
		return err
	}

	incomingDimensionSetKeys := utils.GetMapKeys(incomingDimensionSet)

	var filteredDimensions []map[string]string
	for _, existingDimensionSet := range m.dimensions {
		existingDimensionSetKeys := utils.GetMapKeys(existingDimensionSet)

		if len(existingDimensionSetKeys) != len(incomingDimensionSetKeys) ||
			!utils.AreSlicesEqual(existingDimensionSetKeys, incomingDimensionSetKeys) {
			filteredDimensions = append(filteredDimensions, existingDimensionSet)
		}
	}

	m.dimensions = append(filteredDimensions, incomingDimensionSet)
	return nil

}

func (m *MetricsContext) SetDimensions(dimensionSets []map[string]string, useDefault ...bool) error {
	use := false
	if len(useDefault) > 0 {
		use = useDefault[0]
	}

	for _, dimensionSet := range dimensionSets {
		err := validateDimensionSet(dimensionSet)
		if err != nil {
			return err
		}
	}
	m.shouldUseDefaultDimensions = use
	m.dimensions = dimensionSets
	return nil
}

func (m *MetricsContext) ResetDimensions(useDefault bool) {
	m.shouldUseDefaultDimensions = useDefault
	m.dimensions = make([]map[string]string, 0)
}

func (m *MetricsContext) GetDimensions() []map[string]string {
	if !m.shouldUseDefaultDimensions {
		return m.dimensions
	}

	if len(m.defaultDimensions) == 0 {
		return m.dimensions
	}

	if len(m.dimensions) == 0 {
		return []map[string]string{m.defaultDimensions}
	}

	mergedDimensions := make([]map[string]string, 0)
	for _, customDimension := range m.dimensions {
		m := utils.MergeMaps(m.defaultDimensions, customDimension)
		if m == nil {
			log.Println("Merged dimensions are empty")
		}
		mergedDimensions = append(mergedDimensions, m.(map[string]string))
	}
	return mergedDimensions
}

func (m *MetricsContext) PutMetric(key string, value float64, unit utils.Unit, storageResolution ...utils.StorageResolution) error {
	sR := utils.Standard
	if len(storageResolution) >= 1 {
		sR = storageResolution[0]
	}

	err := validateMetric(key, unit, sR, m.metricNameAndResolutionMap)
	if err != nil {
		return err
	}
	currentMetric := m.Metrics[key]
	if currentMetric.Values != nil && currentMetric.Unit != "" && currentMetric.StorageResolution != 0 {
		currentMetric.addValue(value)
		m.Metrics[key] = currentMetric
	} else {
		m.Metrics[key] = MetricsValue{
			Values:            []float64{value},
			Unit:              unit,
			StorageResolution: sR,
		}
	}
	m.metricNameAndResolutionMap[key] = sR
	return nil
}

func (m *MetricsContext) CreateCopyWithContext(preserveDimensions ...bool) MetricsContext {

	pD := true
	if len(preserveDimensions) > 0 {
		pD = preserveDimensions[0]
	}

	return MetricsContext{
		Namespace:                  m.Namespace,
		Properties:                 m.Properties,
		Metrics:                    make(map[string]MetricsValue),
		Meta:                       map[string]any{"Timestamp": resolveMetaTimestamp(0)},
		dimensions:                 m.dimensions,
		defaultDimensions:          m.defaultDimensions,
		shouldUseDefaultDimensions: pD,
		timestamp:                  m.timestamp,
		metricNameAndResolutionMap: make(map[string]utils.StorageResolution),
	}
}

func resolveMetaTimestamp(timestamp int64) int64 {
	if timestamp == 0 {
		return time.Now().Unix() * 1000
	}
	return timestamp
}

func (m *MetricsContext) Serialize() ([]string, error) {
	var dimensionKeys [][]string
	dimensionProperties := make(map[string]string)

	for _, dimensionSet := range m.GetDimensions() {
		keys := utils.GetMapKeys(dimensionSet)

		if len(keys) > utils.MAX_DIMENSION_SET_SIZE {
			return nil, fmt.Errorf("maximum number of dimensions allowed is %d", utils.MAX_DIMENSION_SET_SIZE)
		}

		dimensionKeys = append(dimensionKeys, keys)
		for k, v := range dimensionSet {
			dimensionProperties[k] = v
		}
	}

	// Function to create the base structure for the JSON object
	createBody := func() map[string]interface{} {
		return map[string]interface{}{
			"_aws": map[string]interface{}{
				"Timestamp": m.Meta["Timestamp"],
				"CloudWatchMetrics": []map[string]interface{}{
					{
						"Namespace":  m.Namespace,
						"Dimensions": dimensionKeys,
						"Metrics":    []interface{}{},
					},
				},
			},
		}
	}

	eventBatches := []string{}
	currentBody := createBody()

	// Function to serialize and add the current body to the batches
	serializeCurrentBody := func() {
		bodyBytes, err := json.Marshal(currentBody)
		if err != nil {
			log.Println("Error serializing JSON:", err)
			return
		}
		eventBatches = append(eventBatches, string(bodyBytes))
		currentBody = createBody()
	}

	// Iterate over the metrics to add them to the event batches
	for key, metric := range m.Metrics {
		// Add metrics in batches, ensuring they don’t exceed limits
		for i := 0; i < len(metric.Values); i += utils.MAX_VALUES_PER_METRIC {
			end := int(math.Min(float64(i+utils.MAX_VALUES_PER_METRIC), float64(len(metric.Values))))
			valueSlice := metric.Values[i:end]

			currentBody[key] = valueSlice
			metricObj := map[string]interface{}{
				"Name":              key,
				"Unit":              string(metric.Unit),
				"StorageResolution": metric.StorageResolution,
			}
			// Only add the metric object if it doesn't already exist
			existingMetrics := currentBody["_aws"].(map[string]interface{})["CloudWatchMetrics"].([]map[string]interface{})[0]["Metrics"].([]interface{})
			metricExists := false
			for _, metric := range existingMetrics {
				if metric.(map[string]interface{})["Name"] == key {
					metricExists = true
					break
				}
			}
			if !metricExists {
				currentBody["_aws"].(map[string]interface{})["CloudWatchMetrics"].([]map[string]interface{})[0]["Metrics"] = append(
					currentBody["_aws"].(map[string]interface{})["CloudWatchMetrics"].([]map[string]interface{})[0]["Metrics"].([]interface{}),
					metricObj,
				)
			}

			// If we hit the max number of metrics per event, serialize the current batch
			if len(currentBody["_aws"].(map[string]interface{})["CloudWatchMetrics"].([]map[string]interface{})[0]["Metrics"].([]interface{})) == utils.MAX_METRICS_PER_EVENT {
				serializeCurrentBody()
			}
		}
	}

	// Serialize the remaining body if there are unprocessed metrics
	if len(currentBody["_aws"].(map[string]interface{})["CloudWatchMetrics"].([]map[string]interface{})[0]["Metrics"].([]interface{})) > 0 {
		serializeCurrentBody()
	}

	return eventBatches, nil
}
