package environments

import (
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/sinks"
)

const (
	Local   = "Local"
	Lambda  = "Lambda"
	Agent   = "Agent"
	EC2     = "EC2"
	ECS     = "ECS"
	Unknown = ""
)

type Environment interface {
	Probe() bool
	GetName() string
	GetType() string
	GetLogGroupName() string
	ConfigureContext(context *context.MetricsContext)
	GetSink() sinks.Sink
}
