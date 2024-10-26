# aws-embedded-metrics
Amazon CloudWatch Embedded Metric Format Client Library for Golang.

<em>I tried to match the standards and structure of the existing client libraries provided by [Amazon Webs Services - Labs](https://github.com/awslabs). Suggestions for improvemnts or bug reports are highly welcome. Just open up an issue or PR.</em>

## Installation

```
go get github.com/tomkalesse/aws-embedded-metrics-go@v0.0.9
```

## Usage

```golang
import "github.com/tomkalesse/aws-embedded-metrics-go/metrics"

metricsLogger := metrics.CreateMetricsLogger()

metricsLogger.PutMetric("TestCount", 1.0, metrics.Count, metrics.StorageResolutionHigh)
metricsLogger.Flush()
```

### ECS

You need to install a CloudWatchAgent while running ECS (Fragate).

### EC2

...

### Lambda

...

### Local

...


## Config

The following keys can be configured as environment variable.
```
"LOG_GROUP_NAME", "LOG_STREAM_NAME", "ENABLE_DEBUG_LOGGING", "SERVICE_NAME", "SERVICE_TYPE", "AGENT_ENDPOINT", "ENVIRONMENT", "NAMESPACE"
```
Optionally the prefix `AWS_EMF_` can be added to all keys.
