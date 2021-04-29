package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSCapacityProviderClinet interface {
	DescribeCapacityProviders(context.Context, *ecs.DescribeCapacityProvidersInput, ...func(*ecs.Options)) (*ecs.DescribeCapacityProvidersOutput, error)
}

func DescribeCapacityProvider(client ECSCapacityProviderClinet, names []string) ([]ecsTypes.CapacityProvider, error) {
	describeCapacityProvidersOutput, err := client.DescribeCapacityProviders(context.TODO(),
		&ecs.DescribeCapacityProvidersInput{
			CapacityProviders: names,
		},
	)
	if err != nil {
		return []ecsTypes.CapacityProvider{}, err
	}
	return describeCapacityProvidersOutput.CapacityProviders, nil
}
