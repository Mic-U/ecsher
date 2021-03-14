package util

import (
	"fmt"
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
		fmt.Println(GetValue(*testIntegerResource))
		t.Fatal("Failed to convert DoubleValue")
	}
}
