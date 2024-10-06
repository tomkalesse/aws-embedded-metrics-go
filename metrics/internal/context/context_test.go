package context

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/utils"
)

func TestCanSetProperty(t *testing.T) {

	context := Empty()
	expectedKey := "key"
	expectedValue := "value"

	context.SetProperty(expectedKey, expectedValue)

	result := context.Properties[expectedKey]

	if result != expectedValue {
		t.Errorf("Expected %s, got %v", expectedValue, result)
	}
}

func TestSetDimensionsAllows30Dimensions(t *testing.T) {

	context := Empty()
	numOfDimensions := 30
	expectedDimensionSet := getDimensionSet(numOfDimensions)

	context.SetDimensions([]map[string]string{expectedDimensionSet})

	if utils.AreMapsEqual(expectedDimensionSet, context.GetDimensions()[0]) != true {
		t.Errorf("Expected %s, got %v", expectedDimensionSet, context.GetDimensions()[0])
	}
}

func TestPutDimensionAddsKeyToDimensionAndSetsDimensionAsProperty(t *testing.T) {

	context := Empty()
	dimension := "testDimension"
	expectedDimension := map[string]string{dimension: dimension}

	context.PutDimensions(expectedDimension)

	if utils.AreMapsEqual(expectedDimension, context.GetDimensions()[0]) != true {
		t.Errorf("Expected %s, got %v", expectedDimension, context.GetDimensions()[0])
	}
}

func TestPutDimensionsAcceptsMultipleUniqueDimensionSets(t *testing.T) {

	context := Empty()
	expectedDimension1 := map[string]string{"d1": "dim1", "d2": "dim2"}
	expectedDimension2 := map[string]string{"d2": "dim2", "d3": "dim3"}

	context.PutDimensions(expectedDimension1)
	context.PutDimensions(expectedDimension2)

	if len(context.GetDimensions()) != 2 {
		t.Errorf("Expected %v, got %v", 2, len(context.GetDimensions()))
	}
	if utils.AreMapsEqual(expectedDimension1, context.GetDimensions()[0]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension1, context.GetDimensions()[0])
	}
	if utils.AreMapsEqual(expectedDimension2, context.GetDimensions()[1]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension2, context.GetDimensions()[1])
	}
}

func TestPutDimensionsWillNotDuplicateDimensions(t *testing.T) {

	context := Empty()
	dimension1 := "d1"
	dimension2 := "d2"
	expectedDimension1 := map[string]string{}
	expectedDimension2 := map[string]string{dimension1: dimension1}
	expectedDimension3 := map[string]string{dimension2: dimension2, dimension1: dimension1}
	expectedDimension4 := map[string]string{dimension2: dimension2}

	context.PutDimensions(map[string]string{})
	context.PutDimensions(map[string]string{dimension1: dimension1})
	context.PutDimensions(map[string]string{dimension1: dimension1, dimension2: dimension2})
	context.PutDimensions(map[string]string{dimension2: dimension2, dimension1: dimension1})
	context.PutDimensions(map[string]string{dimension2: dimension2})
	context.PutDimensions(map[string]string{})
	context.PutDimensions(map[string]string{dimension1: dimension1})
	context.PutDimensions(map[string]string{dimension1: dimension1, dimension2: dimension2})
	context.PutDimensions(map[string]string{dimension2: dimension2, dimension1: dimension1})
	context.PutDimensions(map[string]string{dimension2: dimension2})

	if len(context.GetDimensions()) != 4 {
		t.Errorf("Expected %v, got %v", 4, len(context.GetDimensions()))
	}
	if utils.AreMapsEqual(expectedDimension1, context.GetDimensions()[0]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension1, context.GetDimensions()[0])
	}
	if utils.AreMapsEqual(expectedDimension2, context.GetDimensions()[1]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension2, context.GetDimensions()[1])
	}
	if utils.AreMapsEqual(expectedDimension3, context.GetDimensions()[2]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension3, context.GetDimensions()[2])
	}
	if utils.AreMapsEqual(expectedDimension4, context.GetDimensions()[3]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension4, context.GetDimensions()[3])
	}
}

func TestPutDimensionsWillSortDimensionsCorrectly(t *testing.T) {

	context := Empty()
	dimension1 := "d1"
	dimension2 := "d2"
	expectedDimension1 := map[string]string{dimension2: dimension2, dimension1: dimension1}
	expectedDimension2 := map[string]string{dimension2: dimension2}
	expectedDimension3 := map[string]string{dimension1: dimension1}
	expectedDimension4 := map[string]string{}

	context.PutDimensions(map[string]string{})
	context.PutDimensions(map[string]string{dimension1: dimension1})
	context.PutDimensions(map[string]string{dimension1: dimension1, dimension2: dimension2})
	context.PutDimensions(map[string]string{dimension2: dimension2, dimension1: dimension1})
	context.PutDimensions(map[string]string{dimension2: dimension2})
	context.PutDimensions(map[string]string{dimension1: dimension1, dimension2: dimension2})
	context.PutDimensions(map[string]string{dimension2: dimension2, dimension1: dimension1})
	context.PutDimensions(map[string]string{dimension2: dimension2})
	context.PutDimensions(map[string]string{dimension1: dimension1})
	context.PutDimensions(map[string]string{})

	if len(context.GetDimensions()) != 4 {
		t.Errorf("Expected %v, got %v", 4, len(context.GetDimensions()))
	}
	if utils.AreMapsEqual(expectedDimension1, context.GetDimensions()[0]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension1, context.GetDimensions()[0])
	}
	if utils.AreMapsEqual(expectedDimension2, context.GetDimensions()[1]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension2, context.GetDimensions()[1])
	}
	if utils.AreMapsEqual(expectedDimension3, context.GetDimensions()[2]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension3, context.GetDimensions()[2])
	}
	if utils.AreMapsEqual(expectedDimension4, context.GetDimensions()[3]) != true {
		t.Errorf("Expected %v, got %v", expectedDimension4, context.GetDimensions()[3])
	}
}

func TestGetDimensionsReturnsDefaultDimensionsIfCustomDimensionsNotSet(t *testing.T) {

	context := Empty()
	defaultDimension := "defaultDimension"
	expectedDimension := map[string]string{defaultDimension: defaultDimension}
	context.SetDefaultDimensions(expectedDimension)

	dimensions := context.GetDimensions()

	if utils.AreMapsEqual(expectedDimension, dimensions[0]) != true {
		t.Errorf("Expected %s, got %v", expectedDimension, dimensions[0])
	}
}

func TestGetDimensionsCustomDimensionsPlusDefaultDimensions(t *testing.T) {

	context := Empty()
	defaultDimension := "defaultDimension"
	customDimension := "customDimension"
	expectedDimensions := map[string]string{defaultDimension: defaultDimension, customDimension: customDimension}

	context.SetDefaultDimensions(map[string]string{defaultDimension: defaultDimension})
	context.PutDimensions(map[string]string{customDimension: customDimension})
	dimensions := context.GetDimensions()

	if utils.AreMapsEqual(expectedDimensions, dimensions[0]) != true {
		t.Errorf("Expected %s, got %s", expectedDimensions, dimensions[0])
	}
}

func TestGetDimensionsReturnsOnlyCustomDimensionsIfNoDefaultDimensions(t *testing.T) {

	context := Empty()
	customDimension := "customDimension"
	expectedDimensions := map[string]string{customDimension: customDimension}

	context.PutDimensions(map[string]string{customDimension: customDimension})
	dimensions := context.GetDimensions()

	if utils.AreMapsEqual(expectedDimensions, dimensions[0]) != true {
		t.Errorf("Expected %s, got %v", expectedDimensions, dimensions[0])
	}
}

func TestPutMetricAddsStandardResolutionMetricToMetricsKey(t *testing.T) {

	context := Empty()
	expectedKey := "key"
	expectedValue := 3.14
	expectedUnit := utils.Seconds

	context.PutMetric(expectedKey, expectedValue, expectedUnit, utils.Standard)

	metricDatum := context.Metrics[expectedKey]

	if utils.AreFloat64SlicesEqual([]float64{expectedValue}, metricDatum.Values) != true {
		t.Errorf("Expected %v, got %v", []float64{expectedValue}, metricDatum.Values)
	}
	if expectedUnit != metricDatum.Unit {
		t.Errorf("Expected %v, got %v", expectedUnit, metricDatum.Unit)
	}
}

func TestPutMetricAddsHighResolutionMetricToMetricsKey(t *testing.T) {

	context := Empty()
	expectedKey := "key"
	expectedValue := 3.14
	expectedUnit := utils.Seconds
	expectedStorageResolution := utils.High

	context.PutMetric(expectedKey, expectedValue, expectedUnit, expectedStorageResolution)

	metricDatum := context.Metrics[expectedKey]

	if utils.AreFloat64SlicesEqual([]float64{expectedValue}, metricDatum.Values) != true {
		t.Errorf("Expected %v, got %v", []float64{expectedValue}, metricDatum.Values)
	}
	if expectedUnit != metricDatum.Unit {
		t.Errorf("Expected %v, got %v", expectedUnit, metricDatum.Unit)
	}
	if expectedStorageResolution != metricDatum.StorageResolution {
		t.Errorf("Expected %v, got %v", expectedStorageResolution, metricDatum.StorageResolution)
	}
}

func TestPutMetricUsesNoneUnitIfNotProvided(t *testing.T) {

	context := Empty()
	expectedKey := "key"
	expectedValue := 3.14
	expectedUnit := utils.None

	context.PutMetric(expectedKey, expectedValue, expectedUnit)

	metricDatum := context.Metrics[expectedKey]

	if utils.AreFloat64SlicesEqual([]float64{expectedValue}, metricDatum.Values) != true {
		t.Errorf("Expected %v, got %v", []float64{expectedValue}, metricDatum.Values)
	}
	if expectedUnit != metricDatum.Unit {
		t.Errorf("Expected %v, got %v", expectedUnit, metricDatum.Unit)
	}

}

func TestPutMetricUsesStandardStorageResolutionIfNotProvided(t *testing.T) {

	context := Empty()
	expectedKey := "key"
	expectedValue := 3.14

	context.PutMetric(expectedKey, expectedValue, utils.Seconds)

	metricDatum := context.Metrics[expectedKey]

	if utils.AreFloat64SlicesEqual([]float64{expectedValue}, metricDatum.Values) != true {
		t.Errorf("Expected %v, got %v", []float64{expectedValue}, metricDatum.Values)
	}
	if utils.Standard != metricDatum.StorageResolution {
		t.Errorf("Expected %v, got %v", utils.Standard, metricDatum.StorageResolution)
	}
}

func TestCreateCopyWithContextCopiesPropertiesAndDimensions(t *testing.T) {

	context := Empty()
	context.PutDimensions(map[string]string{"Key": "Value"})
	context.SetProperty("Prop", "Value")

	newContext := context.CreateCopyWithContext()

	if areContextsEqual(newContext, context) != true {
		t.Errorf("Expected %v, got %v", true, false)
	}
}

/*
 *	func TestCreateCopyWithContextCopiesShouldUseDefaultDimensions(t *testing.T) {
 *
 *		context := Empty()
 *		context.SetDimensions([]map[string]string{})
 *		context.SetDefaultDimensions(map[string]string{"Key": "Value"})
 *
 *		newContext := context.CreateCopyWithContext()
 *
 *		if areArrayMapsEqual([]map[string]string{}, newContext.GetDimensions()) != true {
 *			t.Errorf("Expected %v, got %v", []map[string]string{}, newContext.GetDimensions())
 *		}
 *	}
 */

func TestPutDimensionsChecksTheDimensionSetLength(t *testing.T) {

	context := Empty()
	numOfDimensions := 33

	err := context.PutDimensions(getDimensionSet(numOfDimensions))
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestSetDimensionsChecksAllTheDimensionSetsHaveLessThan30Dimensions(t *testing.T) {

	context := Empty()
	numOfDimensions := 33

	err := context.SetDimensions([]map[string]string{getDimensionSet(numOfDimensions)})
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestPutDimensionsWithInvalidDimensionThrowsError(t *testing.T) {
	context := Empty()

	// Test cases for invalid dimensions
	testCases := []struct {
		name          string
		dimensions    map[string]string
		expectedError string
	}{
		{"Invalid Key (Emoji)", map[string]string{"🚀": "value"}, "Invalid dimension"},
		{"Invalid Value (Non-Latin)", map[string]string{"d1": "مارك"}, "Invalid dimension"},
		{"Too Long Key", map[string]string{utils.GenerateString('a', 251): "value"}, "Invalid dimension"},
		{"Too Long Value", map[string]string{"d1": utils.GenerateString('a', 1025)}, "Invalid dimension"},
		{"Empty Key", map[string]string{"": "value"}, "Invalid dimension"},
		{"Empty Value", map[string]string{"d1": ""}, "Invalid dimension"},
		{"Invalid Character in Key", map[string]string{":d1": "value"}, "Invalid dimension"},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := context.PutDimensions(tc.dimensions)
			if err == nil {
				t.Errorf("Expected error but got nil")
			}
		})
	}
}

func TestPutDimensionsWithValidDimensionDoesNotThrowError(t *testing.T) {
	context := Empty()

	// Test cases with valid dimensions
	testCases := []struct {
		name       string
		dimensions map[string]string
	}{
		{"Valid Random Words", map[string]string{utils.GenerateString('a', 5): utils.GenerateString('b', 5)}},
		{"Single Letter Key", map[string]string{"d": utils.GenerateString('c', 5)}},
		{"Valid Word as Value", map[string]string{utils.GenerateString('a', 5): "a"}},
		{"Dimension with Large Value", map[string]string{"dimension:1": utils.GenerateString('a', 1024)}},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := context.PutDimensions(tc.dimensions)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestPutMetricWithInvalidDataThrowsError(t *testing.T) {
	context := Empty()

	var vierfünf utils.StorageResolution = 45
	var null utils.StorageResolution = 0
	var fahrenheit utils.Unit = "Fahrenheit"
	var nilUnit utils.Unit = ""

	// Test cases for invalid metrics
	testCases := []struct {
		metricName       string
		metricValue      float64
		metricUnit       utils.Unit
		metricResolution utils.StorageResolution
		expectedError    string
	}{
		{"", 12345, utils.Count, null, "Invalid metric"},
		{utils.GenerateString('a', utils.MAX_METRIC_NAME_LENGTH+1), 67890, nilUnit, null, "Invalid metric"},
		{"randomWord1", math.MaxFloat64, nilUnit, null, "Invalid metric"},
		{"randomWord2", -math.MaxFloat64, nilUnit, null, "Invalid metric"},
		{"randomWord3", math.MaxInt64 + 1, nilUnit, null, "Invalid metric"},
		{"randomWord4", -math.MaxInt64 - 1, nilUnit, null, "Invalid metric"},
		{"randomWord5", math.NaN(), nilUnit, null, "Invalid metric"},
		{"randomWord6", math.Inf(1), utils.Count, null, "Invalid metric"},
		{"randomWord7", math.Inf(-1), utils.Count, null, "Invalid metric"},
		{"randomWord8", 4, fahrenheit, null, "Invalid metric"},
		{"randomWord9", math.NaN(), utils.Count, null, "Invalid metric"},
		{"randomWordsThree", 123, utils.Seconds, vierfünf, "Invalid metric"},
		{"randomWordsThreeAgain", 456, utils.Seconds, null, "Invalid metric"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("putMetric with name: %s, value: %f, unit: %v and resolution: %v throws error",
			tc.metricName, tc.metricValue, tc.metricUnit, tc.metricResolution), func(t *testing.T) {
			err := context.PutMetric(tc.metricName, tc.metricValue, tc.metricUnit, tc.metricResolution)
			if err == nil {
				t.Errorf("Expected error but got nil")
			}
		})
	}
}

func TestPutMetricWithValidDataDoesNotThrowError(t *testing.T) {

	var null utils.StorageResolution = 0
	var eins utils.StorageResolution = 1

	// Test cases for valid metrics
	testCases := []struct {
		metricName       string
		metricValue      float64
		metricUnit       utils.Unit
		metricResolution utils.StorageResolution
	}{
		{"randomWord1", float64(-1000 - rand.Intn(998)), utils.None, null},
		{"randomWord2", float64(0 - rand.Intn(9999)), utils.Count, null},
		{"randomWordsTwo", float64(0 - rand.Intn(9999)), utils.None, null},
		{"randomWordsThree", float64(0 - rand.Intn(9999)), utils.Seconds, null},
		{"Max_Value", math.MaxInt64, utils.Milliseconds, null},
		{"-Max_Value", -math.MaxInt64, utils.BytesPerSecond, null},
		{"-Max_Value", float64(0 - rand.Intn(9999)), utils.BytesPerSecond, eins},
		{"-Max_Value", float64(0 - rand.Intn(9999)), utils.BytesPerSecond, utils.Standard},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("putMetric with name: %s, value: %f, unit: %v and resolution: %v does not throw error",
			tc.metricName, tc.metricValue, tc.metricUnit, tc.metricResolution), func(t *testing.T) {
			context := Empty()
			err := context.PutMetric(tc.metricName, tc.metricValue, tc.metricUnit, tc.metricResolution)
			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}
		})
	}
}

func TestPutMetricWithSameKeyAndDifferentResolution(t *testing.T) {
	context := Empty()
	expectedKey := "key"
	expectedValue := 3.14
	expectedUnit := utils.Seconds

	context.PutMetric(expectedKey, expectedValue, expectedUnit, utils.Standard)
	err := context.PutMetric(expectedKey, expectedValue, expectedUnit, utils.High)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestSetInvalidNamespaceName(t *testing.T) {

	// Test cases for valid metrics
	testCases := []struct {
		namespace string
	}{
		{""},
		{" "},
		{"àẁş/ćļốṹḓⱳầƭḉⱨ"},
		{"namespace "},
		{utils.GenerateString('a', utils.MAX_NAMESPACE_LENGTH+1)},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Set namspace with name: %s", tc.namespace), func(t *testing.T) {
			context := Empty()
			err := context.SetNamespace(tc.namespace)
			if err == nil {
				t.Errorf("Expected error but got none")
			}
		})
	}
}

func TestSetValidNamespaceName(t *testing.T) {
	context := Empty()
	expectedNamespace := "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz.-_/#:"

	err := context.SetNamespace(expectedNamespace)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestSetInvalidTimestamp(t *testing.T) {

	// Test cases for valid metrics
	testCases := []struct {
		timestamp int64
	}{
		{0},
		{253402300800},
		{time.Now().Unix() - int64(utils.MAX_TIMESTAMP_PAST_AGE) - 1},
		{time.Now().Unix() + int64(utils.MAX_TIMESTAMP_FUTURE_AGE) + 1},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Set timestamp with name: %v", tc.timestamp), func(t *testing.T) {
			context := Empty()
			err := context.SetTimestamp(tc.timestamp)
			if err == nil {
				t.Errorf("Expected error but got none")
			}
		})
	}
}

func TestSetValidTimestamp(t *testing.T) {

	// Test cases for valid metrics
	testCases := []struct {
		timestamp int64
	}{
		{1728121252},
		{time.Now().Add(-utils.MAX_TIMESTAMP_PAST_AGE).Unix() + 1000},
		{time.Now().Add(utils.MAX_TIMESTAMP_FUTURE_AGE).Unix() - 1000},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Set timestamp with name: %v", tc.timestamp), func(t *testing.T) {
			context := Empty()
			err := context.SetTimestamp(tc.timestamp)
			if err != nil {
				t.Errorf("Expected nil but got error %v", err)
			}
		})
	}
}

func getDimensionSet(numOfDimensions int) map[string]string {
	dimensionSet := make(map[string]string)

	for i := 0; i < numOfDimensions; i++ {
		expectedKey := strconv.Itoa(i)
		dimensionSet[expectedKey] = "value" + expectedKey
	}

	return dimensionSet
}

func areContextsEqual(a, b MetricsContext) bool {
	if a.Namespace != b.Namespace {
		return false
	}
	if !utils.AreMapsEqual(a.Properties, b.Properties) {
		return false
	}
	if !areMapsMetricsValueEqual(a.Metrics, b.Metrics) {
		return false
	}
	if !utils.AreMapsAnyEqual(a.Meta, b.Meta) {
		return false
	}
	if !utils.AreArrayMapsEqual(a.dimensions, b.dimensions) {
		return false
	}
	if !utils.AreMapsEqual(a.defaultDimensions, b.defaultDimensions) {
		return false
	}
	if a.shouldUseDefaultDimensions != b.shouldUseDefaultDimensions {
		return false
	}
	if a.timestamp != b.timestamp {
		return false
	}
	if !utils.AreMapsStorageResolutionEqual(a.metricNameAndResolutionMap, b.metricNameAndResolutionMap) {
		return false
	}
	return true
}

func areMapsMetricsValueEqual(map1, map2 map[string]MetricsValue) bool {

	if len(map1) != len(map2) {
		return false
	}

	for key, value := range map1 {
		if v, exists := map2[key]; !exists || !utils.AreFloat64SlicesEqual(value.Values, v.Values) || value.Unit != v.Unit || value.StorageResolution != v.StorageResolution {
			return false
		}
	}

	return true
}
