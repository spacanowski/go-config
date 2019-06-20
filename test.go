package main

import "fmt"

type InternalConfig struct {
	test  string
	test1 string
	test2 string
}

func main() {
	var config InternalConfig

	err := Load(&config)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Printf("Result: %v\n", config)
	fmt.Printf("Result.test: %v\n", config.test)
	fmt.Printf("Result.test1: %v\n", config.test1)
	fmt.Printf("Result.test2: %v\n", config.test2)
}
