package sinks

import "github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"

type Sink interface {
	Name() string
	Accept(context *context.MetricsContext) error
}

type SocketClient interface {
	SendMessage(message []byte) error
}

type Endpoint struct {
	Host     string
	Port     string
	Protocol string
}
