package environments

import (
	"errors"
	"log"
	"sync"

	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/config"
	"github.com/tomkalesse/aws-embedded-metrics-go/metrics/internal/utils"
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

func getEnvironmentFromOverride() (Environment, error) {
	env := config.GetConfig()
	switch env.EnvironmentOverride {
	case utils.Agent:
		return defaultEnvironment, nil
	case utils.EC2:
		return ec2Environment, nil
	case utils.Lambda:
		return lambdaEnvironment, nil
	case utils.ECS:
		return NewECSEnvironment()
	case utils.Local:
		return localEnvironment, nil
	default:
		return nil, errors.New("not found")
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
		var err error
		log.Println("Resolving environment")
		env := config.GetConfig()
		if env.EnvironmentOverride != "" {
			log.Printf("Environment override supplied: %s", env.EnvironmentOverride)
			environment, err = getEnvironmentFromOverride()
			if err == nil {
				return
			}
			log.Printf("Invalid environment provided. Falling back to auto-discovery: %s", env.EnvironmentOverride)
		}
		environment, err = discoverEnvironment()
		if err != nil {
			log.Printf("Failed to discover environment: %v", err)
		}
	})

	if environment == nil {
		return nil, errors.New("failed to resolve environment")
	}
	return environment, nil
}

func CleanResolveEnvironment() (Environment, error) {
	once = sync.Once{}
	return ResolveEnvironment()
}
