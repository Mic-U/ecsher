package ecs

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockCapacityProviderClient struct {
	describeCapacityProvidersOutput ecs.DescribeCapacityProvidersOutput
}

func (m mockCapacityProviderClient) DescribeCapacityProviders(ctx context.Context, params *ecs.DescribeCapacityProvidersInput, optFns ...func(*ecs.Options)) (*ecs.DescribeCapacityProvidersOutput, error) {
	return &m.describeCapacityProvidersOutput, nil
}

func TestDescribeCapacityProvider(t *testing.T) {
	testName := "test"
	client := &mockCapacityProviderClient{
		describeCapacityProvidersOutput: ecs.DescribeCapacityProvidersOutput{
			CapacityProviders: []ecsTypes.CapacityProvider{},
		},
	}

	_, err := DescribeCapacityProvider(client, []string{testName})
	if err != nil {
		t.Fatalf("%d, unexpected error", err)
	}
}
