package util

import (
	"testing"
)

func TestIsValidOutputFormat(t *testing.T) {
	cases := []struct {
		outputFormat string
		expected     bool
	}{
		{
			outputFormat: "default",
			expected:     true,
		},
		{
			outputFormat: "json",
			expected:     true,
		},
		{
			outputFormat: "yaml",
			expected:     true,
		},
	}

	for _, c := range cases {
		actual := IsValidOutputFormat(c.outputFormat)
		if !actual {
			t.Errorf("expect %v, got %v", c.expected, actual)
		}
	}

	actual := IsValidOutputFormat("hoge")
	if actual {
		t.Errorf("expect false, got %v", actual)
	}
}

func TestOutputAsYaml(t *testing.T) {
	type Hoge struct {
		hoge string
		fuga int
	}
	input := Hoge{
		hoge: "abcde",
		fuga: 15,
	}

	_, err := OutputAsYaml(input)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOutputAsJson(t *testing.T) {
	type Hoge struct {
		hoge string
		fuga int
	}
	input := Hoge{
		hoge: "abcde",
		fuga: 15,
	}

	_, err := OutputAsJson(input)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsYamlFormat(t *testing.T) {
	actual1 := IsYamlFormat("yaml")
	if !actual1 {
		t.Errorf("Expected true, but actual false")
	}

	actual2 := IsYamlFormat("yamll")
	if actual2 {
		t.Errorf("Expected false, but actual true")
	}
}

func TestIsJsonFormat(t *testing.T) {
	actual1 := IsJsonFormat("json")
	if !actual1 {
		t.Errorf("Expected true, but actual false")
	}

	actual2 := IsJsonFormat("jsonn")
	if actual2 {
		t.Errorf("Expected false, but actual true")
	}
}

func TestIsDefaultFormat(t *testing.T) {
	actual1 := IsDefaultFormat("default")
	if !actual1 {
		t.Errorf("Expected true, but actual false")
	}

	actual2 := IsDefaultFormat("defaultt")
	if actual2 {
		t.Errorf("Expected false, but actual true")
	}
}
