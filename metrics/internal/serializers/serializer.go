package serializers

import (
	"fmt"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
)

type Serializer interface {
	Serialize(context *context.MetricsContext) []string
}

type FormatVersion string

const (
	VersionZero FormatVersion = "0"
)

func CreateSerializer(version FormatVersion) (Serializer, error) {
	switch version {
	case VersionZero:
		return &LogSerializer{}, nil
	default:
		return nil, fmt.Errorf("unsupported version: %s", version)
	}
}
