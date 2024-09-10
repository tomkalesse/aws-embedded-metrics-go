package sinks

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/serializers"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

var defaultTcpEndpoint = Endpoint{
	Host:     "0.0.0.0",
	Port:     "25888",
	Protocol: TCP,
}

type IEndpoint struct {
	Host     string
	Port     int
	Protocol string
}

type AgentSink struct {
	name          string
	Serializer    serializers.Serializer
	Endpoint      Endpoint
	LogGroupName  string
	LogStreamName string
	SocketClient  SocketClient
	mutex         sync.Mutex
}

func parseEndpoint(endpoint string) Endpoint {
	if endpoint == "" {
		return defaultTcpEndpoint
	}

	parsedURL, err := url.Parse(endpoint)
	if err != nil || parsedURL.Hostname() == "" || parsedURL.Port() == "" || parsedURL.Scheme == "" {
		log.Printf("Failed to parse the provided agent endpoint. Falling back to the default TCP endpoint. %v", err)
		return defaultTcpEndpoint
	}

	port := parsedURL.Port()

	if parsedURL.Scheme != TCP && parsedURL.Scheme != UDP {
		log.Printf("The provided agent endpoint protocol '%s' is not supported. Please use TCP or UDP. Falling back to the default TCP endpoint.", parsedURL.Scheme)
		return defaultTcpEndpoint
	}

	return Endpoint{
		Host:     parsedURL.Hostname(),
		Port:     port,
		Protocol: parsedURL.Scheme,
	}
}

func NewAgentSink(logGroupName, logStreamName string, serializer serializers.Serializer) *AgentSink {
	if serializer == nil {
		serializer = &serializers.LogSerializer{}
	}
	endpoint := parseEndpoint(config.EnvironmentConfig.AgentEndpoint)

	sink := &AgentSink{
		name:          "AgentSink",
		LogGroupName:  logGroupName,
		LogStreamName: logStreamName,
		Serializer:    serializer,
		Endpoint:      endpoint,
		SocketClient:  getSocketClient(endpoint),
	}

	log.Printf("Using socket client: %T", sink.SocketClient)
	return sink
}

func (s *AgentSink) Accept(context *context.MetricsContext) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.LogGroupName != "" {
		context.Meta["LogGroupName"] = s.LogGroupName
	}
	if s.LogStreamName != "" {
		context.Meta["LogStreamName"] = s.LogStreamName
	}

	events := s.Serializer.Serialize(context)
	log.Printf("Sending %d events to socket.", len(events))

	for _, event := range events {
		message := []byte(event + "\n")
		err := s.SocketClient.SendMessage(message)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}

func (s *AgentSink) Name() string {
	return s.name
}

func getSocketClient(endpoint Endpoint) SocketClient {
	log.Printf("Getting socket client for connection: %v", endpoint)
	var client SocketClient
	if endpoint.Protocol == TCP {
		client = NewTcpClient(endpoint)
	} else {
		client = NewUdpClient(endpoint)
	}

	return client
}