package configloader

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

	if err := loadConfigFile(&config, "application-"+profile); err != nil {
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
	for _, extension := range []string{".yaml", ".yml"} {
		if fileName, ok := fileExists(baseFileName + extension); ok {
			return fileName, nil
		}
	}

	return "", fmt.Errorf("No config for file %v with .yaml or .yml extension", baseFileName)
}

func fileExists(fileName string) (string, bool) {
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		log.Printf("No config for file %v\n", fileName)
		return fileName, false
	}

	return fileName, true
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
	valuesMap := config.(map[interface{}]interface{})

	for i := 0; i < res.NumField(); i++ {
		fieldName := res.Type().Field(i).Name

		if _, ok := valuesMap[fieldName]; !ok {
			continue
		}

		field := res.Field(i)
		reflectedField := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
		v, isMap := valuesMap[fieldName].(map[interface{}]interface{})

		if isMap && field.Kind() == reflect.Struct {
			setResult(reflectedField, v)
		} else {
			value := reflect.ValueOf(valuesMap[fieldName])
			if field.Kind() != value.Kind() {
				fmt.Printf("Cannot map. Field '%s' is not of type '%s'\n", fieldName, value.Kind())
				continue
			}
			if value.Kind() == reflect.Slice {
				t := reflect.MakeSlice(field.Type(), value.Len(), value.Cap())
				for i := 0; i < value.Len(); i++ {
					t.Index(i).Set(reflect.ValueOf(value.Index(i).Interface()))
				}
				value = t
			} else if isMap {
				t := reflect.MakeMap(field.Type())
				iter := value.MapRange()
				for iter.Next() {
					t.SetMapIndex(reflect.ValueOf(iter.Key().Interface()), reflect.ValueOf(iter.Value().Interface()))
				}
				value = t
			}
			reflectedField.Set(value)
		}
	}
}
