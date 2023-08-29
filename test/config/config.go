package config

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadYaml(file string) any {
	bytes, _ := os.ReadFile(file)
	var data any
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		panic(err)
	}
	return data
}

func ReadJson(file string) any {
	bytes, _ := os.ReadFile(file)
	var data any
	if err := json.Unmarshal(bytes, &data); err != nil {
		panic(err)
	}
	return data
}
