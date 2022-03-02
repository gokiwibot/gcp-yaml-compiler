package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	// Default file to replace env vars
	pathFile := "./app.yaml"
	// Check if there is a new file path to replace env vars
	if pf := os.Getenv("INPUT_FILE"); pf != "" {
		pathFile = pf
	}
	fmt.Println("Ready to compile ...")
	// Check if file exists
	filename, _ := filepath.Abs(pathFile)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	// Read the file
	var mapResult map[interface{}]interface{}
	err = yaml.Unmarshal(yamlFile, &mapResult)
	if err != nil {
		panic(err)
	}
	fmt.Println("Env variables will be replaced")
	for k, any := range mapResult {
		// Detect the env vars to replace
		if k == "env_variables" {
			err := checkIsPointer(&any)
			if err != nil {
				panic(err)
			}
			valueOf := reflect.ValueOf(any)
			val := reflect.Indirect(valueOf)
			switch val.Type().Kind() {
			case reflect.Map:
				envMap := any.(map[interface{}]interface{})
				for in, iv := range envMap {
					envName := in.(string)
					envVal := iv.(string)
					env := strings.Replace(strings.TrimSpace(envVal), "$", "", -1)
					envMap[envName] = os.Getenv(env)
				}
			default:
				panic(fmt.Sprintf("Error found inside env_variables: %s", val.Type().Kind()))
			}
		}
	}
	fmt.Println("Successfully compiled env variables")
	out, err := yaml.Marshal(mapResult)
	if err != nil {
		panic(err)
	}
	// write the whole body at once
	err = ioutil.WriteFile(pathFile, out, 0644)
	if err != nil {
		panic(err)
	}
}

func checkIsPointer(any interface{}) error {
	if reflect.ValueOf(any).Kind() != reflect.Ptr {
		return fmt.Errorf("you passed something that was not a pointer: %s", any)
	}
	return nil
}
