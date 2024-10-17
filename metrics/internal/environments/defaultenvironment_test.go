package environments

import (
	"os"
	"testing"
)

func TestDefaultEnvironmentProbe(t *testing.T) {

	env := &DefaultEnvironment{}
	result := env.Probe()

	if result != true {
		t.Errorf("Expected true, got %v", result)
	}
}

func TestDefaultEnvironmentGetName(t *testing.T) {

	env := &DefaultEnvironment{}
	result := env.GetName()

	if result != "Unknown" {
		t.Errorf("Expected Unknown, got %v", result)
	}
}

func TestDefaultEnvironmentGetType(t *testing.T) {

	env := &DefaultEnvironment{}
	result := env.GetType()

	if result != "Unknown" {
		t.Errorf("Expected Unknown, got %v", result)
	}
}

func TestDefaultEnvironmentSetName(t *testing.T) {

	expectedValue := "testName"
	os.Setenv("SERVICE_NAME", expectedValue)
	env := &DefaultEnvironment{}
	result := env.GetName()

	if result != expectedValue {
		t.Errorf("Expected %s, got %v", expectedValue, result)
	}
}

func TestDefaultEnvironmentSetType(t *testing.T) {

	expectedValue := "testType"
	os.Setenv("SERVICE_TYPE", expectedValue)
	env := &DefaultEnvironment{}
	result := env.GetType()

	if result != expectedValue {
		t.Errorf("Expected %s, got %v", expectedValue, result)
	}
}

func TestDefaultEnvironmentSetLogGroupName(t *testing.T) {

	expectedValue := "testLogGroup"
	os.Setenv("LOG_GROUP_NAME", expectedValue)
	env := &DefaultEnvironment{}
	result := env.GetLogGroupName()

	if result != expectedValue {
		t.Errorf("Expected %s, got %v", expectedValue, result)
	}
}

func TestDefaultEnvironmentGetLogGroupName(t *testing.T) {

	expectedValue := "testName-metrics"
	serviceName := "testName"
	os.Setenv("LOG_GROUP_NAME", "")
	os.Setenv("SERVICE_NAME", serviceName)
	env := &DefaultEnvironment{}
	result := env.GetLogGroupName()

	if result != expectedValue {
		t.Errorf("Expected %s, got %v", expectedValue, result)
	}
}

func TestDefaultEnvironmentGetSink(t *testing.T) {

	expectedSink := "AgentSink"
	env := &DefaultEnvironment{}
	sink := env.GetSink()

	if sink.Name() != expectedSink {
		t.Errorf("Expected %s, got %v", expectedSink, sink.Name())
	}
}

func TestDefaultEnvironmentGetSinkLogGroupName(t *testing.T) {

	env := &DefaultEnvironment{}
	expectedValue := "testName-metrics"
	serviceName := "testName"
	os.Setenv("LOG_GROUP_NAME", "")
	os.Setenv("SERVICE_NAME", serviceName)
	sink := env.GetSink()

	if sink.LogGroupName() != expectedValue {
		t.Errorf("Expected %s, got %v", expectedValue, sink.LogGroupName())
	}
}
