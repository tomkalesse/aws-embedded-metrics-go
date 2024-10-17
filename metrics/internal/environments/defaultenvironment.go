package environments

import (
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/sinks"
)

type DefaultEnvironment struct {
	sink sinks.Sink
}

func (e *DefaultEnvironment) Probe() bool {
	return true
}

func (e *DefaultEnvironment) GetName() string {
	env := config.GetConfig()
	if env.ServiceName == "" {
		return "Unknown"
	}
	return env.ServiceName
}

func (e *DefaultEnvironment) GetType() string {
	env := config.GetConfig()
	if env.ServiceType == "" {
		return "Unknown"
	}
	return env.ServiceType
}

func (e *DefaultEnvironment) GetLogGroupName() string {
	env := config.GetConfig()
	if env.LogGroupName == "" {
		return e.GetName() + "-metrics"
	}
	return env.LogGroupName
}

func (e *DefaultEnvironment) ConfigureContext(ctx *context.MetricsContext) {
	// no-op
}

func (e *DefaultEnvironment) GetSink() sinks.Sink {
	env := config.GetConfig()
	if e.sink == nil {
		e.sink = sinks.NewAgentSink(e.GetLogGroupName(), env.LogStreamName)
	}
	return e.sink
}
