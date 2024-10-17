package sinks

import (
	"fmt"
	"log"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
)

type ConsoleSink struct {
	name string
}

func NewConsoleSink() *ConsoleSink {
	return &ConsoleSink{
		name: "ConsoleSink",
	}
}

func (s *ConsoleSink) Accept(context *context.MetricsContext) error {
	events, err := context.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize context: %w", err)
	}
	for _, event := range events {
		log.Println(event)
	}
	return nil
}

func (s *ConsoleSink) Name() string {
	return s.name
}

func (s *ConsoleSink) LogGroupName() string {
	return ""
}
