package util

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var validOutputFormats []string = []string{"default", "yaml", "json"}

func IsValidOutputFormat(outputFormat string) bool {
	for _, v := range validOutputFormats {
		if outputFormat == v {
			return true
		}
	}
	return false
}

func OutputAsYaml(input interface{}) (string, error) {
	// TO avoid https://github.com/go-yaml/yaml/issues/463
	jsonInput, _ := OutputAsJson(input)
	yamlInput := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonInput), &yamlInput)
	if err != nil {
		return "", err
	}
	output, err := yaml.Marshal(yamlInput)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func OutputAsArrayedYaml(input interface{}) (string, error) {
	// TO avoid https://github.com/go-yaml/yaml/issues/463
	jsonInput, _ := OutputAsJson(input)
	yamlInput := []interface{}{}
	err := json.Unmarshal([]byte(jsonInput), &yamlInput)
	if err != nil {
		return "", err
	}
	output, err := yaml.Marshal(yamlInput)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func OutputAsJson(input interface{}) (string, error) {
	output, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func IsYamlFormat(outputFormat string) bool {
	return outputFormat == "yaml"
}

func IsJsonFormat(outputFormat string) bool {
	return outputFormat == "json"
}

func IsDefaultFormat(outputFormat string) bool {
	return outputFormat == "default"
}
