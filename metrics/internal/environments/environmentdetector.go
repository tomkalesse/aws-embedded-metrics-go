package environments

import (
	"log"
	"sync"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
)

var (
	lambdaEnvironment  = &LambdaEnvironment{}
	ecsEnvironment     = &ECSEnvironment{}
	ec2Environment     = &EC2Environment{}
	defaultEnvironment = &DefaultEnvironment{}
	localEnvironment   = &LocalEnvironment{}
)

var environments = []Environment{
	lambdaEnvironment,
	ecsEnvironment,
	ec2Environment,
}

var environment Environment
var once sync.Once

func getEnvironmentFromOverride() Environment {
	switch config.EnvironmentConfig.EnvironmentOverride {
	case Agent:
		return defaultEnvironment
	case EC2:
		return ec2Environment
	case Lambda:
		return lambdaEnvironment
	case ECS:
		return ecsEnvironment
	case Local:
		return localEnvironment
	default:
		return nil
	}
}

func discoverEnvironment() (Environment, error) {
	log.Println("Discovering environment")
	for _, env := range environments {
		log.Printf("Testing: %T", env)

		if err := env.Probe(); err {
			return env, nil
		} else {
			log.Printf("Failed probe: %T", env)
		}
	}
	return defaultEnvironment, nil
}

func ResolveEnvironment() (Environment, error) {
	once.Do(func() {
		log.Println("Resolving environment")
		if config.EnvironmentConfig.EnvironmentOverride != "" {
			log.Printf("Environment override supplied: %s", config.EnvironmentConfig.EnvironmentOverride)
			environment = getEnvironmentFromOverride()
			if environment != nil {
				return
			}
			log.Printf("Invalid environment provided. Falling back to auto-discovery: %s", config.EnvironmentConfig.EnvironmentOverride)
		}

		var err error
		environment, err = discoverEnvironment()
		if err != nil {
			log.Printf("Failed to discover environment: %v", err)
		}
	})

	return environment, nil
}

func CleanResolveEnvironment() (Environment, error) {
	once = sync.Once{}
	return ResolveEnvironment()
}
