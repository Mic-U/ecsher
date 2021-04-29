package util

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func TestGetCapacityProviderType(t *testing.T) {
	cases := []struct {
		cp       ecsTypes.CapacityProvider
		expected string
	}{
		{
			cp: ecsTypes.CapacityProvider{
				Name: aws.String("FARGATE"),
			},
			expected: "FARGATE",
		},
		{
			cp: ecsTypes.CapacityProvider{
				Name: aws.String("FARGATE_SPOT"),
			},
			expected: "FARGATE",
		},
		{
			cp: ecsTypes.CapacityProvider{
				Name: aws.String("TEST"),
			},
			expected: "ASG",
		},
	}

	for _, c := range cases {
		actual := GetCapacityProviderType(c.cp)
		if actual != c.expected {
			t.Errorf("expect %v, got %v", c.expected, actual)
		}
	}
}
