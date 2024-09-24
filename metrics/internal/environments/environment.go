package environments

import (
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/sinks"
)

type Environment interface {
	Probe() bool
	GetName() string
	GetType() string
	GetLogGroupName() string
	ConfigureContext(context *context.MetricsContext)
	GetSink() sinks.Sink
}
