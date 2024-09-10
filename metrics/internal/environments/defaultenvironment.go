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
	if config.EnvironmentConfig.ServiceName == "" {
		return "Unknown"
	}
	return config.EnvironmentConfig.ServiceName
}

func (e *DefaultEnvironment) GetType() string {
	if config.EnvironmentConfig.ServiceType == "" {
		return "Unknown"
	}
	return config.EnvironmentConfig.ServiceType
}

func (e *DefaultEnvironment) GetLogGroupName() string {
	if config.EnvironmentConfig.LogGroupName == "" {
		return ""
	}
	return e.GetName() + "-metrics"
}

func (e *DefaultEnvironment) ConfigureContext(ctx *context.MetricsContext) {
	// no-op
}

func (e *DefaultEnvironment) GetSink() sinks.Sink {
	if e.sink == nil {
		e.sink = sinks.NewAgentSink(e.GetLogGroupName(), config.EnvironmentConfig.LogStreamName, nil)
	}
	return e.sink
}
