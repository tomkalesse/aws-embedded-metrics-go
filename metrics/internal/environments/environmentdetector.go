package environments

import (
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

func getEnvironmentFromOverride() Environment {
	env := config.GetConfig()
	switch env.EnvironmentOverride {
	case utils.Agent:
		return defaultEnvironment
	case utils.EC2:
		return ec2Environment
	case utils.Lambda:
		return lambdaEnvironment
	case utils.ECS:
		return ecsEnvironment
	case utils.Local:
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
		env := config.GetConfig()
		if env.EnvironmentOverride != "" {
			log.Printf("Environment override supplied: %s", env.EnvironmentOverride)
			environment = getEnvironmentFromOverride()
			if environment != nil {
				return
			}
			log.Printf("Invalid environment provided. Falling back to auto-discovery: %s", env.EnvironmentOverride)
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
