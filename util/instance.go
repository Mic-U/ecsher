package util

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

const (
	IntegerType = "INTEGER"
	LongType    = "LONG"
	DoubleType  = "DOUBLE"
)

func GetRemainingCpuString(remainingResources []types.Resource) string {
	for _, remainingResource := range remainingResources {
		if *remainingResource.Name == "CPU" {
			return GetValue(remainingResource)
		}
	}
	return "0"
}

func GetRemainingMemoryString(remainingResources []types.Resource) string {
	for _, remainingResource := range remainingResources {
		if *remainingResource.Name == "MEMORY" {
			return GetValue(remainingResource)
		}
	}
	return "0"
}

func GetValue(remainingResource types.Resource) string {
	switch *remainingResource.Type {
	case IntegerType:
		return strconv.Itoa(int(remainingResource.IntegerValue))

	case LongType:
		return strconv.Itoa(int(remainingResource.LongValue))

	case DoubleType:
		return strconv.FormatFloat(remainingResource.DoubleValue, 'f', 2, 64)

	default:
		return "0"
	}
}
