package serializers

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/utils"
)

type LogSerializer struct{}

type MetricProgress struct {
	Name    string
	NumLeft int
}

func (s *LogSerializer) Serialize(context *context.MetricsContext) []string {
	var dimensionKeys [][]string
	dimensionProperties := make(map[string]string)

	for _, dimensionSet := range context.GetDimensions() {
		keys := make([]string, 0, len(dimensionSet))
		for k := range dimensionSet {
			keys = append(keys, k)
		}

		if len(keys) > utils.MAX_DIMENSION_SET_SIZE {
			errMsg := fmt.Sprintf(
				"Maximum number of dimensions allowed are %d. Account for default dimensions if not using set_dimensions.",
				utils.MAX_DIMENSION_SET_SIZE)
			panic(errors.New(errMsg))
		}

		dimensionKeys = append(dimensionKeys, keys)
		for k, v := range dimensionSet {
			dimensionProperties[k] = v
		}
	}

	createBody := func() map[string]interface{} {
		return map[string]interface{}{
			"_aws": map[string]interface{}{
				"CloudWatchMetrics": []map[string]interface{}{
					{
						"Dimensions": dimensionKeys,
						"Metrics":    []map[string]interface{}{},
						"Namespace":  context.Namespace,
					},
				},
			},
			"Meta":       context.Meta,
			"Properties": context.Properties,
			"Dimensions": dimensionProperties,
		}
	}

	eventBatches := []string{}
	currentBody := createBody()

	currentMetricsInBody := func() int {
		return len(currentBody["_aws"].(map[string]interface{})["CloudWatchMetrics"].([]map[string]interface{})[0]["Metrics"].([]map[string]interface{}))
	}

	hasMaxMetrics := func() bool {
		return currentMetricsInBody() == utils.MAX_METRICS_PER_EVENT
	}

	serializeCurrentBody := func() {
		body, _ := json.Marshal(currentBody)
		eventBatches = append(eventBatches, string(body))
		currentBody = createBody()
	}

	remainingMetrics := maxHeapify(context.Metrics)
	processedMetrics := []MetricProgress{}

	for !remainingMetrics.isEmpty() {
		metricProgress := remainingMetrics.extractRoot()
		metric := context.Metrics[metricProgress.Name]
		startIndex := len(metric.Values) - metricProgress.NumLeft
		var metricValue interface{}
		if metricProgress.NumLeft == 1 {
			metricValue = metric.Values[startIndex]
		} else {
			endIndex := startIndex + utils.MAX_VALUES_PER_METRIC
			if endIndex > len(metric.Values) {
				endIndex = len(metric.Values)
			}
			metricValue = metric.Values[startIndex:endIndex]
		}

		currentBody[metricProgress.Name] = metricValue

		metricBody := map[string]interface{}{
			"Name": metricProgress.Name,
			"Unit": metric.Unit,
		}
		if metric.StorageResolution == utils.High {
			metricBody["StorageResolution"] = utils.High
		}

		currentBody["_aws"].(map[string]interface{})["CloudWatchMetrics"].([]map[string]interface{})[0]["Metrics"] = append(
			currentBody["_aws"].(map[string]interface{})["CloudWatchMetrics"].([]map[string]interface{})[0]["Metrics"].([]map[string]interface{}), metricBody)

		metricProgress.NumLeft -= utils.MAX_VALUES_PER_METRIC
		if metricProgress.NumLeft > 0 {
			processedMetrics = append(processedMetrics, metricProgress)
		}

		if hasMaxMetrics() || remainingMetrics.isEmpty() {
			serializeCurrentBody()
			for _, m := range processedMetrics {
				remainingMetrics.insert(m)
			}
			processedMetrics = []MetricProgress{}
		}
	}

	if len(eventBatches) == 0 || currentMetricsInBody() > 0 {
		serializeCurrentBody()
	}

	return eventBatches
}

func maxHeapify(metrics map[string]context.MetricsValue) *MaxHeap {
	heap := &MaxHeap{}
	for name, value := range metrics {
		heap.insert(MetricProgress{
			Name:    name,
			NumLeft: len(value.Values),
		})
	}
	heap.init(func(a, b MetricProgress) bool {
		return a.NumLeft > b.NumLeft
	})
	return heap
}

type MaxHeap struct {
	items []MetricProgress
}

func (h *MaxHeap) insert(item MetricProgress) {
	h.items = append(h.items, item)
	sort.SliceStable(h.items, func(i, j int) bool {
		return h.items[i].NumLeft > h.items[j].NumLeft
	})
}

func (h *MaxHeap) extractRoot() MetricProgress {
	if len(h.items) == 0 {
		return MetricProgress{}
	}
	root := h.items[0]
	h.items = h.items[1:]
	return root
}

func (h *MaxHeap) isEmpty() bool {
	return len(h.items) == 0
}

func (h *MaxHeap) init(compare func(a, b MetricProgress) bool) {
	// Assuming the MaxHeap Init method
	// compare function should be passed here
}
