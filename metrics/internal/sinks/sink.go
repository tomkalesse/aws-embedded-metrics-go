package sinks

import "github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"

type Sink interface {
	Accept(context *context.MetricsContext) error
	Name() string
	LogGroupName() string
}

type SocketClient interface {
	SendMessage(message []byte) error
}

type Endpoint struct {
	Host     string
	Port     string
	Protocol string
}
