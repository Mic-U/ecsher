package util

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func TestGetValue(t *testing.T) {
	integerType := "INTEGER"
	testIntegerResource := &types.Resource{
		DoubleValue:  23.14,
		IntegerValue: 3,
		LongValue:    4,
		Type:         &integerType,
	}
	if GetValue(*testIntegerResource) != "3" {
		t.Fatal("Failed to convert IntegerValue")
	}

	longType := "LONG"
	testIntegerResource.Type = &longType
	if GetValue(*testIntegerResource) != "4" {
		t.Fatal("Failed to convert LongValue")
	}

	doubleType := "DOUBLE"
	testIntegerResource.Type = &doubleType
	if GetValue(*testIntegerResource) != "23.14" {
		t.Fatal("Failed to convert DoubleValue")
	}

	otherType := "OTHER"
	testIntegerResource.Type = &otherType
	if GetValue(*testIntegerResource) != "0" {
		t.Fatal("Failed to convert DoubleValue")
	}
}

func TestGetRemainingCpuString(t *testing.T) {
	CPU := "CPU"
	memory := "MEMORY"
	integer := "INTEGER"
	cases := []struct {
		remainingResources []types.Resource
		expected           string
	}{
		{
			remainingResources: []types.Resource{
				{
					Name:         &CPU,
					Type:         &integer,
					IntegerValue: 10,
				},
				{
					Name:         &memory,
					Type:         &integer,
					IntegerValue: 20,
				},
			},
			expected: "10",
		},
		{
			remainingResources: []types.Resource{
				{
					Name:         &memory,
					Type:         &integer,
					IntegerValue: 20,
				},
			},
			expected: "0",
		},
	}

	for _, c := range cases {
		actual := GetRemainingCpuString(c.remainingResources)
		if actual != c.expected {
			t.Errorf("expect %v, got %v", c.expected, actual)
		}
	}
}

func TestGetRemainingMemoryString(t *testing.T) {
	CPU := "CPU"
	memory := "MEMORY"
	integer := "INTEGER"
	cases := []struct {
		remainingResources []types.Resource
		expected           string
	}{
		{
			remainingResources: []types.Resource{
				{
					Name:         &CPU,
					Type:         &integer,
					IntegerValue: 10,
				},
				{
					Name:         &memory,
					Type:         &integer,
					IntegerValue: 20,
				},
			},
			expected: "20",
		},
		{
			remainingResources: []types.Resource{
				{
					Name:         &CPU,
					Type:         &integer,
					IntegerValue: 20,
				},
			},
			expected: "0",
		},
	}

	for _, c := range cases {
		actual := GetRemainingMemoryString(c.remainingResources)
		if actual != c.expected {
			t.Errorf("expect %v, got %v", c.expected, actual)
		}
	}
}
