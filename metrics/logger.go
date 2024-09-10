package metrics

import (
	"log"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/environments"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/utils"
)

type MetricsLogger struct {
	context                 context.MetricsContext
	environment             environments.Environment
	flushPreserveDimensions bool
}

func CreateMetricsLogger() MetricsLogger {
	context := context.Empty()
	environment, err := environments.ResolveEnvironment()
	if err != nil {
		log.Printf("Error resolving environment: %v", err)
	}
	return MetricsLogger{context, environment, true}
}

func (l *MetricsLogger) Flush() {
	environment, err := environments.ResolveEnvironment()
	if err != nil {
		log.Printf("Error resolving environment: %v", err)
	}
	l.configureContextForEnvironment(&l.context, environment)
	sink := environment.GetSink()
	sink.Accept(&l.context)
	l.context = l.context.CreateCopyWithContext(l.flushPreserveDimensions)
}

func (l *MetricsLogger) SetProperty(key string, value string) MetricsLogger {
	l.context.SetProperty(key, value)
	return *l
}

func (l *MetricsLogger) PutDimensions(dimensions map[string]string) MetricsLogger {
	l.context.PutDimensions(dimensions)
	return *l
}

func (l *MetricsLogger) SetDimensions(dimensionSetOrSets interface{}, useDefault ...bool) MetricsLogger {
	defaultValue := false
	if len(useDefault) > 0 {
		defaultValue = useDefault[0]
	}

	switch v := dimensionSetOrSets.(type) {
	case []map[string]string:
		l.context.SetDimensions(v, defaultValue)
	case map[string]string:
		l.context.SetDimensions([]map[string]string{v}, defaultValue)
	default:
		log.Println("Invalid type for dimensionSetOrSets")
	}
	return *l
}

func (l *MetricsLogger) ResetDimensions(useDefault bool) MetricsLogger {
	l.context.ResetDimensions(useDefault)
	return *l
}

func (l *MetricsLogger) PutMetric(key string, value float64, unit utils.Unit, storageResolution utils.StorageResolution) MetricsLogger {
	l.context.PutMetric(key, value, unit, storageResolution)
	return *l
}

func (l *MetricsLogger) SetNamespace(value string) MetricsLogger {
	l.context.SetNamespace(value)
	return *l
}

func (l *MetricsLogger) SetTimestamp(value int64) MetricsLogger {
	l.context.SetTimestamp(value)
	return *l
}

func (l *MetricsLogger) New() *MetricsLogger {
	m := context.MetricsContext{}
	environment, err := environments.ResolveEnvironment()
	if err != nil {
		log.Printf("Error resolving environment: %v", err)
	}
	return &MetricsLogger{m.CreateCopyWithContext(true), environment, true}
}

func (l *MetricsLogger) configureContextForEnvironment(context *context.MetricsContext, environment environments.Environment) {

	serviceName := config.EnvironmentConfig.ServiceName
	if serviceName == "" {
		serviceName = environment.GetName()
	}
	serviceType := config.EnvironmentConfig.ServiceType
	if serviceType == "" {
		serviceType = environment.GetType()
	}

	defaultDimensions := map[string]string{
		"LogGroup":    environment.GetLogGroupName(),
		"ServiceName": serviceName,
		"ServiceType": serviceType,
	}
	context.SetDefaultDimensions(defaultDimensions)
	environment.ConfigureContext(context)
}