package configloader

import (
	"os"
	"testing"
)

type innerStruct struct {
	innerField string
}

type internalConfig struct {
	test      string
	internal1 struct {
		test1     int
		internal2 struct {
			test2 bool
			test3 string
			empty string
			test4 []int
			test5 []string
			test6 map[string]int
			test7 []innerStruct
			test8 map[string]innerStruct
		}
	}
}

type internalConfigWrongStructType struct {
	test      string
	internal1 string
}

type internalConfigWrongFieldType struct {
	test int
}

func TestLoadProperties(t *testing.T) {
	t.Run("-p=dev", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"cmd", "-p=dev"}
		fullTest(t)
	})

	t.Run("--profile=dev", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"cmd", "--profile=dev"}
		fullTest(t)
	})
}

func TestLoadPropertiesWithoutProfile(t *testing.T) {
	setup()
	defer teardown()

	var config internalConfig

	if err := Load(&config); err != nil {
		t.Fatalf("Properties loading failed %v", err)
	}

	if config.test != "ok" {
		t.Fatalf("test failed, expected: ok  actual: %v", config.test)
	}

	if config.internal1.test1 != 0 {
		t.Fatalf("test1 failed, expected: 0  actual: %v", config.internal1.test1)
	}

	if config.internal1.internal2.test2 != false {
		t.Fatalf("test2 failed, expected: false  actual: %v", config.internal1.internal2.test2)
	}

	if config.internal1.internal2.test3 != "" {
		t.Fatalf("test3 failed, expected: ''  actual: %v", config.internal1.internal2.test3)
	}

	if config.internal1.internal2.empty != "" {
		t.Fatalf("test4 failed, expected: ''  actual: %v", config.internal1.internal2.empty)
	}
}

func TestShouldOmitWrongStructType(t *testing.T) {
	setup()
	defer teardown()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "--profile=dev"}

	var config internalConfigWrongStructType

	if err := Load(&config); err != nil {
		t.Fatalf("Properties loading failed %v", err)
	}

	if config.test != "dev-ok" {
		t.Fatalf("test failed, expected: ok  actual: %v", config.test)
	}

	if config.internal1 != "" {
		t.Fatalf("internal1 failed, expected: ''  actual: %v", config.internal1)
	}
}

func TestShouldOmitWrongFieldType(t *testing.T) {
	setup()
	defer teardown()

	var config internalConfigWrongFieldType

	if err := Load(&config); err != nil {
		t.Fatalf("Properties loading failed %v", err)
	}

	if config.test != 0 {
		t.Fatalf("test failed, expected: 0  actual: %v", config.test)
	}
}

func fullTest(t *testing.T) {
	setup()
	defer teardown()

	var config internalConfig

	if err := Load(&config); err != nil {
		t.Fatalf("Properties loading failed %v", err)
	}

	if config.test != "dev-ok" {
		t.Fatalf("test failed, expected: dev-ok  actual: %v", config.test)
	}

	if config.internal1.test1 != 1 {
		t.Fatalf("test1 failed, expected: 1  actual: %v", config.internal1.test1)
	}

	if config.internal1.internal2.test2 != true {
		t.Fatalf("test2 failed, expected: true  actual: %v", config.internal1.internal2.test2)
	}

	if config.internal1.internal2.test3 != "test4-ok" {
		t.Fatalf("test3 failed, expected: test4-ok  actual: %v", config.internal1.internal2.test3)
	}

	if len(config.internal1.internal2.test4) != 3 {
		t.Fatalf("test4 failed, expected: 3 actual: %v", len(config.internal1.internal2.test4))
	}

	if config.internal1.internal2.test4[1] != 2 {
		t.Fatalf("test4 failed, expected: 2 actual: %v", config.internal1.internal2.test4[1])
	}

	if len(config.internal1.internal2.test5) != 2 {
		t.Fatalf("test5 failed, expected: 2 actual: %v", len(config.internal1.internal2.test5))
	}

	if config.internal1.internal2.test5[1] != "dsa" {
		t.Fatalf("test5 failed, expected: 'dsa' actual: %v", config.internal1.internal2.test5[1])
	}

	if len(config.internal1.internal2.test6) != 2 {
		t.Fatalf("test6 failed, expected: 2 actual: %v", len(config.internal1.internal2.test6))
	}

	if config.internal1.internal2.test6["key1"] != 6 {
		t.Fatalf("test6 failed, expected: 6 actual: %v", config.internal1.internal2.test6["key1"])
	}

	if len(config.internal1.internal2.test7) != 2 {
		t.Fatalf("test7 failed, expected: 2 actual: %v", len(config.internal1.internal2.test7))
	}

	if config.internal1.internal2.test7[0].innerField != "in1" {
		t.Fatalf("test7 failed, expected: 'in1' actual: %v", config.internal1.internal2.test7[0].innerField)
	}

	if config.internal1.internal2.test7[1].innerField != "in2" {
		t.Fatalf("test7 failed, expected: 'in2' actual: %v", config.internal1.internal2.test7[1].innerField)
	}

	if len(config.internal1.internal2.test8) != 2 {
		t.Fatalf("test8 failed, expected: 2 actual: %v", len(config.internal1.internal2.test8))
	}

	if config.internal1.internal2.test8["k1"].innerField != "i1" {
		t.Fatalf("test8 failed, expected: 'i1' actual: %v", config.internal1.internal2.test8["k1"].innerField)
	}

	if config.internal1.internal2.test8["k2"].innerField != "i2" {
		t.Fatalf("test8 failed, expected: 'i2' actual: %v", config.internal1.internal2.test8["k2"].innerField)
	}

	if config.internal1.internal2.empty != "" {
		t.Fatalf("test9 failed, expected: ''  actual: %v", config.internal1.internal2.empty)
	}
}

const applicationYamlFileName string = "application.yaml"
const applicationDevYamlFileName string = "application-dev.yml"
const applicationYaml string = `test: ok
test11: app`
const applicationDevYaml string = `test: dev-ok
test2: ok
internal1:
  test1: 1
  internal2:
    test2: true
    test3: test4-ok
    test4: [1,2,3]
    test5:
      - asd
      - dsa
    test6:
      key1: 6
      key2: 7
    test7:
      - innerField: in1
      - innerField: in2
    test8:
      k1:
        innerField: i1
      k2:
        innerField: i2`

func setup() {
	createTestFile(applicationYamlFileName, applicationYaml)
	createTestFile(applicationDevYamlFileName, applicationDevYaml)
}

func createTestFile(fileName string, fileContent string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if f, err := os.Create(fileName); err == nil {
			f.WriteString(fileContent)
		}
	}
}

func teardown() {
	deleteTestFile(applicationYamlFileName)
	deleteTestFile(applicationDevYamlFileName)
}

func deleteTestFile(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
	}
}
