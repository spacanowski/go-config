package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"unsafe"

	"gopkg.in/yaml.v2"
)

// Load fills config from file
func Load(config interface{}) error {
	res := reflect.ValueOf(config).Elem()

	if config == nil {
		return errors.New("No config specified")
	}

	if err := loadConfigFile(&config, "application"); err != nil {
		return err
	}

	setResult(res, config)

	profile := getProfile()

	if profile == "" {
		return nil
	}

	if err := loadConfigFile(&config, fmt.Sprintf("application-%s", profile)); err != nil {
		return err
	}

	setResult(res, config)

	return nil
}

func loadConfigFile(config interface{}, baseFileName string) error {
	fileName, err := getPropertyFileName(baseFileName)
	if err != nil {
		return err
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

func getPropertyFileName(baseFileName string) (string, error) {
	fileName := baseFileName + ".yaml"
	if _, err := os.Stat(fileName); err == nil {
		return fileName, nil
	}

	fileName = baseFileName + ".yml"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Printf("No config for file %v\n", fileName)

		return "", fmt.Errorf("No config for file %v", fileName)
	}

	return fileName, nil
}

func getProfile() string {
	var result string

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if strings.HasPrefix(arg, "--profile=") || strings.HasPrefix(arg, "-p=") {
			result = strings.SplitAfter(arg, "=")[1]

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
