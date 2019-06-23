package configloader

import (
	"os"
	"testing"
)

type internalConfig struct {
	test      string
	internal1 struct {
		test1     int
		internal2 struct {
			test2 bool
			test3 string
			empty string
		}
	}
}

type internalConfigWrong struct {
	test      string
	internal1 string
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

func TestShouldOmitWrongStructNameClash(t *testing.T) {
	setup()
	defer teardown()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "--profile=dev"}

	var config internalConfigWrong

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

	if config.internal1.internal2.empty != "" {
		t.Fatalf("test4 failed, expected: ''  actual: %v", config.internal1.internal2.empty)
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
    test3: test4-ok`

func setup() {
	createTestFile(applicationYamlFileName, applicationYaml)
	createTestFile(applicationDevYamlFileName, applicationDevYaml)
}

func createTestFile(fileName string, fileContent string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if f, cerr := os.Create(fileName); cerr == nil {
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
