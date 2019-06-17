package main

import "fmt"

type InternalConfig struct {
	test string
}

func main() {
	var config InternalConfig

	err := Load(&config)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("Result: %v\n", config)
	fmt.Printf("Result.test: %v\n", config.test)
}
