package metrics

import (
	"log"
	"log/slog"
	"os"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/environments"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/utils"
)

var slogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type MetricsLogger struct {
	context                 context.MetricsContext
	environment             environments.Environment
	flushPreserveDimensions bool
}

func CreateMetricsLogger() MetricsLogger {
	ctx := context.Empty()
	environment, err := environments.ResolveEnvironment()
	if err != nil {
		log.Println("Error resolving environment: " + err.Error())
	}
	return MetricsLogger{ctx, environment, true}
}

func (l *MetricsLogger) Flush() {
	environment, err := environments.ResolveEnvironment()
	if err != nil {
		msg := "Error resolving environment: " + err.Error()
		slogger.Error(msg)
	}
	l.configureContextForEnvironment(&l.context, environment)
	sink := environment.GetSink()
	sink.Accept(&l.context)
	l.context = l.context.CreateCopyWithContext(l.flushPreserveDimensions)
}

func (l *MetricsLogger) SetProperty(key string, value string) {
	l.context.SetProperty(key, value)
}

func (l *MetricsLogger) PutDimensions(dimensions map[string]string) {
	l.context.PutDimensions(dimensions)
}

func (l *MetricsLogger) SetDimensions(dimensionSetOrSets interface{}, useDefault ...bool) {
	defaultValue := false
	if len(useDefault) > 0 {
		defaultValue = useDefault[0]
	}

	switch v := dimensionSetOrSets.(type) {
	case []map[string]string:
		err := l.context.SetDimensions(v, defaultValue)
		if err != nil {
			slogger.Error(err.Error())
		}
	case map[string]string:
		err := l.context.SetDimensions([]map[string]string{v}, defaultValue)
		if err != nil {
			slogger.Error(err.Error())
		}
	default:
		log.Println("Invalid type for dimensionSetOrSets")
	}
}

func (l *MetricsLogger) ResetDimensions(useDefault bool) {
	l.context.ResetDimensions(useDefault)
}

func (l *MetricsLogger) PutMetric(key string, value float64, unit utils.Unit, storageResolution utils.StorageResolution) {
	err := l.context.PutMetric(key, value, unit, storageResolution)
	if err != nil {
		slogger.Error(err.Error())
	}
}

func (l *MetricsLogger) SetNamespace(value string) {
	l.context.SetNamespace(value)
}

func (l *MetricsLogger) SetTimestamp(value int64) {
	l.context.SetTimestamp(value)
}

func (l *MetricsLogger) New() *MetricsLogger {
	environment, err := environments.ResolveEnvironment()
	if err != nil {
		log.Println("Error resolving environment: " + err.Error())
	}
	return &MetricsLogger{l.context.CreateCopyWithContext(true), environment, true}
}

func (l *MetricsLogger) configureContextForEnvironment(context *context.MetricsContext, environment environments.Environment) {
	env := config.GetConfig()
	serviceName := env.ServiceName
	if serviceName == "" {
		serviceName = environment.GetName()
	}
	serviceType := env.ServiceType
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
