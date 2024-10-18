package metrics

import (
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {

	os.Setenv("ENVIRONMENT", "ECS")

	logger := CreateMetricsLogger()
	if 2 == 3 {
		t.Errorf("Failed to create logger")
	}

	logger.PutMetric("test", 1.0, Count, StorageResolutionStandard)
	logger.PutMetric("test", 2.0, Count, StorageResolutionStandard)
	logger.PutMetric("test1", 1.0, Count, StorageResolutionStandard)
	logger.Flush()
}
