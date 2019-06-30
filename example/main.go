package main

import (
	"fmt"

	configloader "https://github.com/spacanowski/go-config/loader"
)

type myAppConfig struct {
	appName     string
	accountName string
	aws         struct {
		accountNumber int
		regions       []string
		clientSqs     map[string]string
		db            struct {
			autoCommit bool
			url        string
			username   string
		}
	}
}

func main() {
	var config myAppConfig
	if err := configloader.Load(&config); err != nil {
		// Handle config load error
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("app name: %v\n", config.appName)
	fmt.Printf("account name: %v\n", config.accountName)
	fmt.Printf("aws account number: %v\n", config.aws.accountNumber)
	fmt.Printf("aws regions: %v\n", config.aws.regions)
	fmt.Printf("aws clientSqs: %v\n", config.aws.clientSqs)
	fmt.Printf("aws db autoCommit: %v\n", config.aws.db.autoCommit)
	fmt.Printf("aws db url: %v\n", config.aws.db.url)
	fmt.Printf("aws db username: %v\n", config.aws.db.username)
}
