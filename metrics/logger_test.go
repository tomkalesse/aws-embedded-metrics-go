package metrics

import (
	"testing"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
)

func TestIntegration(t *testing.T) {

	config.EnvironmentConfig.EnvironmentOverride = "Local"

	logger := CreateMetricsLogger()
	if 2 == 3 {
		t.Errorf("Failed to create logger")
	}

	logger.PutMetric("test", 1.0, Count, StorageResolutionStandard)
	logger.Flush()
}
