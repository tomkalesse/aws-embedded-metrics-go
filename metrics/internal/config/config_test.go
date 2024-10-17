package config

import (
	"os"
	"testing"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/utils"
)

func TestSetLogGroupName(t *testing.T) {

	var expectedValue = "testLogGroup"
	os.Setenv("AWS_EMF_LOG_GROUP_NAME", expectedValue)
	env := GetConfig()
	if env.LogGroupName != expectedValue {
		t.Errorf("Failed to set log group name, expected %s, got %s", expectedValue, env.LogGroupName)
	}
	if env.LogGroupName != expectedValue {
		t.Errorf("Failed to set log group name. expected %s, got %s", expectedValue, env.LogGroupName)
	}

}

func TestSetLogStreamName(t *testing.T) {

	var expectedValue = "testLogStream"
	os.Setenv("AWS_EMF_LOG_STREAM_NAME", expectedValue)
	env := GetConfig()
	if env.LogStreamName != expectedValue {
		t.Errorf("Failed to set log stream name, expected %s, got %s", expectedValue, env.LogStreamName)
	}
	if env.LogStreamName != expectedValue {
		t.Errorf("Failed to set log stream name, expected %s, got %s", expectedValue, env.LogStreamName)
	}

}

func TestEnableDebugLogging(t *testing.T) {

	var expectedValue = true
	os.Setenv("AWS_EMF_ENABLE_DEBUG_LOGGING", "true")
	env := GetConfig()
	if env.DebuggingLoggingEnabled != expectedValue {
		t.Errorf("Failed to enable debug logging")
	}
	if env.DebuggingLoggingEnabled != expectedValue {
		t.Errorf("Failed to enable debug logging")
	}

}

func TestSetServiceName(t *testing.T) {

	var expectedValue = "testService"
	os.Setenv("AWS_EMF_SERVICE_NAME", expectedValue)
	env := GetConfig()
	if env.ServiceName != expectedValue {
		t.Errorf("Failed to set service name, expected %s, got %s", expectedValue, env.ServiceName)
	}
	if env.ServiceName != expectedValue {
		t.Errorf("Failed to set service name, expected %s, got %s", expectedValue, env.ServiceName)
	}

}

func TestSetServiceNameShort(t *testing.T) {

	var expectedValue = "testService"
	os.Setenv("SERVICE_NAME", expectedValue)
	env := GetConfig()
	if env.ServiceName != expectedValue {
		t.Errorf("Failed to set service name, expected %s, got %s", expectedValue, env.ServiceName)
	}
	if env.ServiceName != expectedValue {
		t.Errorf("Failed to set service name, expected %s, got %s", expectedValue, env.ServiceName)
	}

}

func TestSetServiceNamePrecedence(t *testing.T) {

	var expectedValue1 = "testService"
	var expectedValue2 = "testServiceWithPrefix"
	os.Setenv("SERVICE_NAME", expectedValue1)
	os.Setenv("AWS_EMF_SERVICE_NAME", expectedValue2)
	env := GetConfig()
	if env.ServiceName != expectedValue2 {
		t.Errorf("Failed to set service name, expected %s, got %s", expectedValue2, env.ServiceName)
	}
	if env.ServiceName != expectedValue2 {
		t.Errorf("Failed to set service name, expected %s, got %s", expectedValue2, env.ServiceName)
	}

}

func TestSetServiceType(t *testing.T) {

	var expectedValue = "testServiceType"
	os.Setenv("AWS_EMF_SERVICE_TYPE", expectedValue)
	env := GetConfig()
	if env.ServiceType != expectedValue {
		t.Errorf("Failed to ser service type, expected %s, got %s", expectedValue, env.ServiceType)
	}
	if env.ServiceType != expectedValue {
		t.Errorf("Failed to ser service type, expected %s, got %s", expectedValue, env.ServiceType)
	}

}

func TestSetServiceTypeShort(t *testing.T) {

	var expectedValue = "testServiceType"
	os.Setenv("SERVICE_TYPE", expectedValue)
	env := GetConfig()
	if env.ServiceType != expectedValue {
		t.Errorf("Failed to ser service type, expected %s, got %s", expectedValue, env.ServiceType)
	}
	if env.ServiceType != expectedValue {
		t.Errorf("Failed to ser service type, expected %s, got %s", expectedValue, env.ServiceType)
	}

}

func TestSetServiceTypePrecedence(t *testing.T) {

	var expectedValue1 = "testServiceType"
	var expectedValue2 = "testServiceTypeWithPrefix"
	os.Setenv("SERVICE_TYPE", expectedValue1)
	os.Setenv("AWS_EMF_SERVICE_TYPE", expectedValue2)
	env := GetConfig()
	if env.ServiceType != expectedValue2 {
		t.Errorf("Failed to ser service type, expected %s, got %s", expectedValue2, env.ServiceType)
	}
	if env.ServiceType != expectedValue2 {
		t.Errorf("Failed to ser service type, expected %s, got %s", expectedValue2, env.ServiceType)
	}

}

func TestSetAgentEndpoint(t *testing.T) {

	var expectedValue = "https://testEndpoint:1234"
	os.Setenv("AWS_EMF_AGENT_ENDPOINT", expectedValue)
	env := GetConfig()
	if env.AgentEndpoint != expectedValue {
		t.Errorf("Failed to set agent endpoint, expected %s, got %s", expectedValue, env.AgentEndpoint)
	}
	if env.AgentEndpoint != expectedValue {
		t.Errorf("Failed to set agent endpoint, expected %s, got %s", expectedValue, env.AgentEndpoint)
	}

}

func TestSetEnvironment(t *testing.T) {

	var expectedValue = "Local"
	os.Setenv("AWS_EMF_ENVIRONMENT", expectedValue)
	env := GetConfig()
	if env.EnvironmentOverride != utils.Local {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.EnvironmentOverride)
	}
	if env.EnvironmentOverride != utils.Local {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.EnvironmentOverride)
	}

}

func TestSetEnvironmentDefault(t *testing.T) {

	var expectedValue = ""
	os.Setenv("AWS_EMF_ENVIRONMENT", expectedValue)
	env := GetConfig()
	if env.EnvironmentOverride != utils.Unknown {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.EnvironmentOverride)
	}
	if env.EnvironmentOverride != utils.Unknown {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.EnvironmentOverride)
	}

}

func TestSetEnvironmentRandom(t *testing.T) {

	var expectedValue = "notExistingEnvironment"
	os.Setenv("AWS_EMF_ENVIRONMENT", expectedValue)
	env := GetConfig()
	if env.EnvironmentOverride != utils.Unknown {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.EnvironmentOverride)
	}
	if env.EnvironmentOverride != utils.Unknown {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.EnvironmentOverride)
	}

}

func TestDefaultNamespace(t *testing.T) {

	var expectedValue = "aws-embedded-metrics"
	env := GetConfig()
	if env.Namespace != expectedValue {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.Namespace)
	}
	if env.Namespace != expectedValue {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.Namespace)
	}

}

func TestSetNamespace(t *testing.T) {

	var expectedValue = "namespace"
	os.Setenv("AWS_EMF_NAMESPACE", expectedValue)
	env := GetConfig()
	if env.Namespace != expectedValue {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.Namespace)
	}
	if env.Namespace != expectedValue {
		t.Errorf("Failed to set environment, expected %s, got %s", expectedValue, env.Namespace)
	}

}
