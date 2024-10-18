package sinks

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
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
	Endpoint      Endpoint
	logGroupName  string
	logStreamName string
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

func NewAgentSink(logGroupName, logStreamName string) *AgentSink {
	env := config.GetConfig()
	endpoint := parseEndpoint(env.AgentEndpoint)
	sink := &AgentSink{
		name:          "AgentSink",
		logGroupName:  logGroupName,
		logStreamName: logStreamName,
		Endpoint:      endpoint,
		SocketClient:  getSocketClient(endpoint),
	}
	log.Printf("Using socket client: %T", sink.SocketClient)
	return sink
}

func (s *AgentSink) Accept(context *context.MetricsContext) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.logGroupName != "" {
		context.Meta["LogGroupName"] = s.logGroupName
	}
	if s.logStreamName != "" {
		context.Meta["LogStreamName"] = s.logStreamName
	}

	events, err := context.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize context: %w", err)
	}
	for _, event := range events {
		message := []byte(event + "\n")
		log.Printf("Sending message: %s", message)
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

func (s *AgentSink) LogGroupName() string {
	return s.logGroupName
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
