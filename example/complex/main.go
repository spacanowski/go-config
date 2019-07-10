package main

import (
	"fmt"

	configloader "github.com/spacanowski/go-config/loader"
)

type feature struct {
	desc      string
	featureID int
	enabled   bool
}

type user struct {
	id    int
	login string
}

type internalConfig struct {
	appName  string
	features map[string]feature
	aws      struct {
		regions []string
		users   []user
		db      struct {
			url      string
			username string
		}
		sqs struct {
			clientQ map[string]string
		}
	}
}

func main() {
	var config internalConfig
	if err := configloader.Load(&config); err != nil {
		// Handle config load error
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("app name: %v\n", config.appName)
	fmt.Printf("features: %v\n", config.features)
	fmt.Printf("aws regions: %v\n", config.aws.regions)
	fmt.Printf("aws users: %v\n", config.aws.users)
	fmt.Printf("aws db: %v\n", config.aws.db)
	fmt.Printf("aws queues: %v\n", config.aws.sqs.clientQ)
}
