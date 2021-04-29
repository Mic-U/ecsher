package util

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

const (
	FARGATE_TYPE = "FARGATE"
	ASG_TYPE     = "ASG"
)

func GetCapacityProviderType(capacityProvider ecsTypes.CapacityProvider) string {
	if *capacityProvider.Name == "FARGATE" || *capacityProvider.Name == "FARGATE_SPOT" {
		return FARGATE_TYPE
	} else {
		return ASG_TYPE
	}
}
