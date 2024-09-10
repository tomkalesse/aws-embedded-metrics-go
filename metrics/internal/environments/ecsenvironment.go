package environments

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/sinks"
)

type ECSMetadataResponse struct {
	Name               string               `json:"Name"`
	DockerId           string               `json:"DockerId"`
	DockerName         string               `json:"DockerName"`
	Image              string               `json:"Image"`
	FormattedImageName string               `json:"FormattedImageName"`
	ImageID            string               `json:"ImageID"`
	Ports              string               `json:"Ports"`
	Labels             ECSMetadataLabels    `json:"Labels"`
	CreatedAt          string               `json:"CreatedAt"`
	StartedAt          string               `json:"StartedAt"`
	Networks           []ECSMetadataNetwork `json:"Networks"`
}

type ECSMetadataLabels struct {
	Cluster               string `json:"com.amazonaws.ecs.cluster"`
	ContainerName         string `json:"com.amazonaws.ecs.container-name"`
	TaskArn               string `json:"com.amazonaws.ecs.task-arn"`
	TaskDefinitionFamily  string `json:"com.amazonaws.ecs.task-definition-family"`
	TaskDefinitionVersion string `json:"com.amazonaws.ecs.task-definition-version"`
}

type ECSMetadataNetwork struct {
	NetworkMode   string   `json:"NetworkMode"`
	IPv4Addresses []string `json:"IPv4Addresses"`
}

type ECSEnvironment struct {
	sink              sinks.Sink
	metadata          *ECSMetadataResponse
	fluentBitEndpoint string
	mutex             sync.Mutex
}

func formatImageName(imageName string) string {
	if imageName == "" {
		return imageName
	}
	parts := strings.Split(imageName, "/")
	return parts[len(parts)-1]
}

func (e *ECSEnvironment) Probe() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	ecsMetadataURI := os.Getenv("ECS_CONTAINER_METADATA_URI")
	if ecsMetadataURI == "" {
		return false
	}

	fluentHost := os.Getenv("FLUENT_HOST")
	if fluentHost != "" && config.EnvironmentConfig.AgentEndpoint == "" {
		e.fluentBitEndpoint = fmt.Sprintf("tcp://%s:%d", fluentHost, 25888)
		config.EnvironmentConfig.AgentEndpoint = e.fluentBitEndpoint
		fmt.Printf("Using FluentBit configuration. Endpoint: %s\n", e.fluentBitEndpoint)
	}

	u, err := url.Parse(ecsMetadataURI)
	if err != nil {
		log.Println("Failed to parse ECS_CONTAINER_METADATA_URI:", err)
		return false
	}

	resp, err := http.Get(u.String())
	if err != nil {
		log.Println("Failed to collect ECS Container Metadata:", err)
		return false
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&e.metadata); err != nil {
		log.Println("Error decoding ECS metadata:", err)
		return false
	}

	if e.metadata != nil {
		e.metadata.FormattedImageName = formatImageName(e.metadata.Image)
		log.Println("Successfully collected ECS Container metadata.")
	}

	return true
}

func (e *ECSEnvironment) GetName() string {
	if config.EnvironmentConfig.ServiceName != "" {
		return config.EnvironmentConfig.ServiceName
	}

	if e.metadata != nil && e.metadata.FormattedImageName != "" {
		return e.metadata.FormattedImageName
	}
	return "Unknown"
}

func (e *ECSEnvironment) GetType() string {
	return "AWS::ECS::Container"
}

func (e *ECSEnvironment) GetLogGroupName() string {
	if e.fluentBitEndpoint != "" {
		return ""
	}

	if config.EnvironmentConfig.LogGroupName != "" {
		return config.EnvironmentConfig.LogGroupName
	}

	return e.GetName()
}

func (e *ECSEnvironment) ConfigureContext(ctx *context.MetricsContext) {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	e.addProperty(ctx, "containerId", hostname)
	e.addProperty(ctx, "createdAt", e.metadata.CreatedAt)
	e.addProperty(ctx, "startedAt", e.metadata.StartedAt)
	e.addProperty(ctx, "image", e.metadata.Image)
	e.addProperty(ctx, "cluster", e.metadata.Labels.Cluster)
	e.addProperty(ctx, "taskArn", e.metadata.Labels.TaskArn)

	if e.fluentBitEndpoint != "" {
		ctx.SetDefaultDimensions(map[string]string{
			"ServiceName": config.EnvironmentConfig.ServiceName,
			"ServiceType": e.GetType(),
		})
	}
}

func (e *ECSEnvironment) GetSink() sinks.Sink {
	if e.sink == nil {
		e.sink = sinks.NewAgentSink(e.GetLogGroupName(), config.EnvironmentConfig.LogStreamName, nil)
	}
	return e.sink
}

func (e *ECSEnvironment) addProperty(ctx *context.MetricsContext, key, value string) {
	if value != "" {
		ctx.SetProperty(key, value)
	}
}
