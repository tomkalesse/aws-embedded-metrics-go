package sinks

import (
	"log"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/serializers"
)

type ConsoleSink struct {
	name       string
	Serializer serializers.Serializer
}

func NewConsoleSink(serializer serializers.Serializer) *ConsoleSink {
	if serializer == nil {
		serializer = &serializers.LogSerializer{}
	}

	return &ConsoleSink{
		name:       "ConsoleSink",
		Serializer: serializer,
	}
}

func (s *ConsoleSink) Accept(context *context.MetricsContext) error {
	events := s.Serializer.Serialize(context)
	for _, event := range events {
		log.Println(event)
	}
	return nil
}

func (s *ConsoleSink) Name() string {
	return s.name
}
