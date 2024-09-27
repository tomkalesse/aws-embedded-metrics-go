# aws-embedded-metrics
Amazon CloudWatch Embedded Metric Format Client Library for Golang.

I tried to match the standards and structure of the existing client libraries provided by [Amazon Webs Services - Labs](https://github.com/awslabs). Suggestions for improvemnts or bug reports are highly welcome. Just open up an issue.

Compared to other Go libraries for embedded metrics, this implementation supports environments such as EC2, ECS, and Lambda, utilizing sinks for metric flushing directly to a CloudWatch agent.
