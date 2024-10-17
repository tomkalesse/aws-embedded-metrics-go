package environments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/context"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/sinks"
)

const host = "169.254.169.254"
const tokenPath = "/latest/api/token"
const tokenRequestHeaderKey = "X-aws-ec2-metadata-token-ttl-seconds"
const tokenRequestHeaderValue = "21600"
const metadataPath = "/latest/dynamic/instance-identity/document"
const metadataRequestTokenHeaderKey = "X-aws-ec2-metadata-token"

type EC2MetadataResponse struct {
	ImageId          string `json:"imageId"`
	AvailabilityZone string `json:"availabilityZone"`
	PrivateIp        string `json:"privateIp"`
	InstanceId       string `json:"instanceId"`
	InstanceType     string `json:"instanceType"`
}

type EC2Environment struct {
	metadata *EC2MetadataResponse
	sink     sinks.Sink
	token    string
	mutex    sync.Mutex
}

// Probe fetches the token and EC2 metadata to determine if the environment is an EC2 instance.
func (e *EC2Environment) Probe() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Fetch token
	token, err := e.fetchToken()
	if err != nil {
		log.Println("Error fetching token:", err)
		return false
	}
	e.token = token

	// Fetch metadata
	metadata, err := e.fetchMetadata(token)
	if err != nil {
		log.Println("Error fetching metadata:", err)
		return false
	}
	e.metadata = metadata

	return e.metadata != nil
}

// Fetch token from EC2 metadata service
func (e *EC2Environment) fetchToken() (string, error) {
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s%s", host, tokenPath), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set(tokenRequestHeaderKey, tokenRequestHeaderValue)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokenBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

// Fetch EC2 instance metadata
func (e *EC2Environment) fetchMetadata(token string) (*EC2MetadataResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s%s", host, metadataPath), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(metadataRequestTokenHeaderKey, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metadata EC2MetadataResponse
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

// GetName returns the service name or "Unknown" if not configured.
func (e *EC2Environment) GetName() string {
	env := config.GetConfig()
	if env.ServiceName == "" {
		log.Println("Unknown ServiceName.")
		return "Unknown"
	}
	return env.ServiceName
}

// GetType returns the environment type, which is 'AWS::EC2::Instance' if metadata is available.
func (e *EC2Environment) GetType() string {
	if e.metadata != nil {
		return "AWS::EC2::Instance"
	}
	return "Unknown"
}

// GetLogGroupName returns the log group name.
func (e *EC2Environment) GetLogGroupName() string {
	env := config.GetConfig()
	if env.LogGroupName != "" {
		return env.LogGroupName
	}
	return fmt.Sprintf("%s-metrics", e.GetName())
}

// ConfigureContext adds EC2 metadata to the provided MetricsContext.
func (e *EC2Environment) ConfigureContext(ctx *context.MetricsContext) {
	if e.metadata != nil {
		ctx.SetProperty("imageId", e.metadata.ImageId)
		ctx.SetProperty("instanceId", e.metadata.InstanceId)
		ctx.SetProperty("instanceType", e.metadata.InstanceType)
		ctx.SetProperty("privateIP", e.metadata.PrivateIp)
		ctx.SetProperty("availabilityZone", e.metadata.AvailabilityZone)
	}
}

// GetSink returns the sink for the EC2 environment.
func (e *EC2Environment) GetSink() sinks.Sink {
	env := config.GetConfig()
	if e.sink == nil {
		e.sink = sinks.NewAgentSink(e.GetLogGroupName(), env.LogStreamName)
	}
	return e.sink
}
