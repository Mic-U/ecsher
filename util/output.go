package util

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
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
	output, err := yaml.Marshal(input)
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
