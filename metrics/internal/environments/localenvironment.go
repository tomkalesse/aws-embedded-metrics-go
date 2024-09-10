package environments

import (
	"log"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/sinks"
)

type LocalEnvironment struct {
	sink sinks.Sink
}

func (e *LocalEnvironment) Probe() bool {
	return false
}

func (e *LocalEnvironment) GetName() string {
	if config.EnvironmentConfig.ServiceName == "" {
		log.Println("Unknown ServiceName.")
		return "Unknown"
	}
	return config.EnvironmentConfig.ServiceName
}

func (e *LocalEnvironment) GetType() string {
	if config.EnvironmentConfig.ServiceType == "" {
		log.Println("Unknown ServiceType.")
		return "Unknown"
	}
	return config.EnvironmentConfig.ServiceType
}

func (e *LocalEnvironment) GetLogGroupName() string {
	if config.EnvironmentConfig.LogGroupName != "" {
		return config.EnvironmentConfig.LogGroupName
	}
	return e.GetName() + "-metrics"
}

func (e *LocalEnvironment) ConfigureContext(ctx *context.MetricsContext) {
	// no-op
}

func (e *LocalEnvironment) GetSink() sinks.Sink {
	if e.sink == nil {
		e.sink = &sinks.ConsoleSink{}
	}
	return e.sink
}
