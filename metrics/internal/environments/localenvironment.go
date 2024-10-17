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
	env := config.GetConfig()
	if env.ServiceName == "" {
		log.Println("Unknown ServiceName.")
		return "Unknown"
	}
	return env.ServiceName
}

func (e *LocalEnvironment) GetType() string {
	env := config.GetConfig()
	if env.ServiceType == "" {
		log.Println("Unknown ServiceType.")
		return "Unknown"
	}
	return env.ServiceType
}

func (e *LocalEnvironment) GetLogGroupName() string {
	env := config.GetConfig()
	if env.LogGroupName != "" {
		return env.LogGroupName
	}
	return e.GetName() + "-metrics"
}

func (e *LocalEnvironment) ConfigureContext(ctx *context.MetricsContext) {
	// no-op
}

func (e *LocalEnvironment) GetSink() sinks.Sink {
	if e.sink == nil {
		e.sink = sinks.NewConsoleSink()
	}
	return e.sink
}
