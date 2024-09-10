package config

import "os"

const (
	ENV_VAR_PREFIX             = "AWS_EMF"
	MAX_DIMENSION_SET_SIZE     = 30
	MAX_DIMENSION_NAME_LENGTH  = 250
	MAX_DIMENSION_VALUE_LENGTH = 1024
	MAX_METRIC_NAME_LENGTH     = 1024
	MAX_NAMESPACE_LENGTH       = 256
	VALID_NAMESPACE_REGEX      = "^[a-zA-Z0-9._#:/-]+$"
	VALID_DIMENSION_REGEX      = "^[\x00-\x7F]+$"
	MAX_TIMESTAMP_PAST_AGE     = 1209600000 // 2 weeks
	MAX_TIMESTAMP_FUTURE_AGE   = 7200000    // 2 hours
	DEFAULT_NAMESPACE          = "aws-embedded-metrics"
	MAX_METRICS_PER_EVENT      = 100
	MAX_VALUES_PER_METRIC      = 100
	DEFAULT_AGENT_HOST         = "0.0.0.0"
	DEFAULT_AGENT_PORT         = 25888
)

type configKeys struct {
	LOG_GROUP_NAME       string
	LOG_STREAM_NAME      string
	ENABLE_DEBUG_LOGGING string
	SERVICE_NAME         string
	SERVICE_TYPE         string
	AGENT_ENDPOINT       string
	ENVIRONMENT_OVERRIDE string
	NAMESPACE            string
}

var ConfigKeys = configKeys{
	LOG_GROUP_NAME:       "LOG_GROUP_NAME",
	LOG_STREAM_NAME:      "LOG_STREAM_NAME",
	ENABLE_DEBUG_LOGGING: "ENABLE_DEBUG_LOGGING",
	SERVICE_NAME:         "SERVICE_NAME",
	SERVICE_TYPE:         "SERVICE_TYPE",
	AGENT_ENDPOINT:       "AGENT_ENDPOINT",
	ENVIRONMENT_OVERRIDE: "ENVIRONMENT_OVERRIDE",
	NAMESPACE:            "NAMESPACE",
}

type Config struct {
	DebuggingLoggingEnabled bool
	ServiceName             string
	ServiceType             string
	LogGroupName            string
	LogStreamName           string
	AgentEndpoint           string
	EnvironmentOverride     string
	Namespace               string
}

var EnvironmentConfig = Config{
	DebuggingLoggingEnabled: tryGetEnvVariableAsBoolean(ConfigKeys.ENABLE_DEBUG_LOGGING, false),
	ServiceName:             getEnvVar(ConfigKeys.SERVICE_NAME),
	ServiceType:             getEnvVar(ConfigKeys.SERVICE_TYPE),
	LogGroupName:            getEnvVar(ConfigKeys.LOG_GROUP_NAME),
	LogStreamName:           getEnvVar(ConfigKeys.LOG_STREAM_NAME),
	AgentEndpoint:           getEnvVar(ConfigKeys.AGENT_ENDPOINT),
	EnvironmentOverride:     getEnvVar(ConfigKeys.ENVIRONMENT_OVERRIDE),
	Namespace:               getNamespace(ConfigKeys.NAMESPACE),
}

func getEnvVar(key string) string {
	if os.Getenv(ENV_VAR_PREFIX+"_"+key) == "" {
		return os.Getenv(key)
	}
	return os.Getenv(ENV_VAR_PREFIX + "_" + key)
}

func getNamespace(key string) string {
	if os.Getenv(key) == "" {
		return DEFAULT_NAMESPACE
	}
	return os.Getenv(key)
}

func tryGetEnvVariableAsBoolean(key string, fallback bool) bool {
	value := getEnvVar(key)
	if value == "" {
		return fallback
	}
	return value == "true"
}
