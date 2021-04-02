package util

import "testing"

func TestExtractCommand(t *testing.T) {
	testCase := []string{"tast", "ls", "-al"}
	result := ExtractCommand(testCase)
	expected := "ls -al"
	if result != expected {
		t.Errorf("expected is '%s', but actual is '%s'", expected, result)
	}
}
