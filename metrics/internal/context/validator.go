package context

import (
	"errors"
	"strconv"
	"strings"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/utils"
)

func validateNamespace(namespace string) error {
	if namespace == "" {
		return errors.New("namespace cannot be empty")
	}
	if len(namespace) > utils.MAX_NAMESPACE_LENGTH {
		return errors.New("namespace cannot be longer than " + strconv.Itoa(utils.MAX_NAMESPACE_LENGTH) + " characters")
	}
	if !utils.VALID_NAMESPACE_REGEX.MatchString(namespace) {
		return errors.New("namespace must match the pattern " + utils.VALID_NAMESPACE_REGEX.String())
	}
	return nil
}

func validateDimensionSet(dimensionSet map[string]string) error {

	if len(dimensionSet) > utils.MAX_DIMENSION_SET_SIZE {
		return errors.New("maximum number of dimensions per dimension set allowed are " + strconv.Itoa(utils.MAX_DIMENSION_SET_SIZE))
	}

	for key, value := range dimensionSet {
		value = strings.TrimSpace(value)
		dimensionSet[key] = value

		if !utils.VALID_DIMENSION_REGEX.MatchString(key) {
			return errors.New("dimension key " + key + " has invalid characters ")
		}
		if !utils.VALID_DIMENSION_REGEX.MatchString(value) {
			return errors.New("Dimension value " + value + " has invalid characters")
		}
		if strings.TrimSpace(key) == "" {
			return errors.New("dimension key " + key + " must include at least one non-whitespace character")
		}
		if value == "" {
			return errors.New("dimension value " + value + " must include at least one non-whitespace character")
		}
		if len(key) > utils.MAX_DIMENSION_NAME_LENGTH {
			return errors.New("dimension key " + key + " must not exceed maximum length")
		}
		if len(value) > utils.MAX_DIMENSION_VALUE_LENGTH {
			return errors.New("dimension value " + value + " must not exceed maximum length")
		}
		if strings.HasPrefix(key, ":") {
			return errors.New("dimension key " + key + " cannot start with ':'")
		}
	}

	return nil
}

func validateMetric(key string, unit utils.Unit, storageResolution utils.StorageResolution, metricNameAndResolutionMap map[string]utils.StorageResolution) error {

	if len(strings.TrimSpace(key)) == 0 {
		return errors.New("metric key " + key + "must include at least one non-whitespace character")
	}
	if len(key) > utils.MAX_METRIC_NAME_LENGTH {
		return errors.New("metric key " + key + " must not exceed maximum length")
	}
	if unit == "" || !isValidUnit(unit) {
		return errors.New("metric unit " + string(unit) + " is not a valid")
	}
	if storageResolution == 0 || !isValidStorageResolution(storageResolution) {
		return errors.New("metric storage resolution " + strconv.Itoa(int(storageResolution)) + " is not a valid")
	}
	if metricNameAndResolutionMap[key] != 0 && metricNameAndResolutionMap[key] != storageResolution {
		return errors.New("resolution for metrics " + key + " is already set. A single log event cannot have a metric with two different resolutions.")
	}

	return nil
}

func isValidUnit(unit utils.Unit) bool {
	for _, u := range utils.Units {
		if u == unit {
			return true
		}
	}
	return false
}

func isValidStorageResolution(resolution utils.StorageResolution) bool {
	for _, r := range utils.StorageResolutions {
		if r == resolution {
			return true
		}
	}
	return false
}
