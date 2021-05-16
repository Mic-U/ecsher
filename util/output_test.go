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
