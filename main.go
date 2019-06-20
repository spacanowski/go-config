package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"unsafe"

	"gopkg.in/yaml.v2"
)

// Load fills config from file
func Load(config interface{}) error {
	res := reflect.ValueOf(config).Elem()

	if config == nil {
		return errors.New("No config specified")
	}

	if err := loadConfigFile(&config, "application.yaml"); err != nil {
		return err
	}

	setResult(res, config)

	profile := getProfile()

	if profile == "" {
		return nil
	}

	if err := loadConfigFile(&config, fmt.Sprintf("application-%s.yaml", profile)); err != nil {
		return err
	}

	setResult(res, config)

	return nil
}

func loadConfigFile(config interface{}, fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Printf("No config for file %v\n", fileName)

		return nil
	}

	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Cannot load config for %v: #%v\n", fileName, err)

		return fmt.Errorf("Cannot load config for %v: #%v", fileName, err)
	}

	log.Printf("Loading config from %v\n", fileName)

	if err := yaml.Unmarshal(configFile, config); err != nil {
		log.Fatalf("Failed to map config: %v\n", err)

		return fmt.Errorf("Failed to map config %v", err)
	}

	return nil
}

func getProfile() string {
	var result string

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if arg == "--profile" || arg == "-p" {
			if len(os.Args) >= i+1 {
				result = os.Args[i+1]
			}

			break
		}
	}

	return result
}

func setResult(res reflect.Value, config interface{}) {
	for i := 0; i < res.NumField(); i++ {
		f := res.Field(i)
		rf := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
		m := config.(map[interface{}]interface{})
		k := res.Type().Field(i).Name
		if _, ok := m[k]; !ok {
			continue
		}
		rf.Set(reflect.ValueOf(m[k]))
	}
}
