package environments

import (
	"os"
	"strings"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/sinks"
)

type LambdaEnvironment struct {
	sink sinks.Sink
}

func (e *LambdaEnvironment) Probe() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

func (e *LambdaEnvironment) GetName() string {
	name := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if name != "" {
		return name
	}
	return "Unknown"
}

func (e *LambdaEnvironment) GetType() string {
	return "AWS::Lambda::Function"
}

func (e *LambdaEnvironment) GetLogGroupName() string {
	return e.GetName()
}

func (e *LambdaEnvironment) ConfigureContext(ctx *context.MetricsContext) {
	e.addProperty(ctx, "executionEnvironment", os.Getenv("AWS_EXECUTION_ENV"))
	e.addProperty(ctx, "memorySize", os.Getenv("AWS_LAMBDA_FUNCTION_MEMORY_SIZE"))
	e.addProperty(ctx, "functionVersion", os.Getenv("AWS_LAMBDA_FUNCTION_VERSION"))
	e.addProperty(ctx, "logStreamId", os.Getenv("AWS_LAMBDA_LOG_STREAM_NAME"))

	if trace := e.getSampledTrace(); trace != "" {
		e.addProperty(ctx, "traceId", trace)
	}
}

func (e *LambdaEnvironment) GetSink() sinks.Sink {
	if e.sink == nil {
		e.sink = sinks.NewConsoleSink()
	}
	return e.sink
}

func (e *LambdaEnvironment) addProperty(ctx *context.MetricsContext, key, value string) {
	if value != "" {
		ctx.SetProperty(key, value)
	}
}

func (e *LambdaEnvironment) getSampledTrace() string {
	traceID := os.Getenv("_X_AMZN_TRACE_ID")
	if traceID != "" && strings.Contains(traceID, "Sampled=1") {
		return traceID
	}
	return ""
}
