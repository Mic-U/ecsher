package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func GetService(cluster string) ([]ecsTypes.Service, error) {
	client := GetClient()
	listServicesOutput, err := client.ListServices(context.TODO(),
		&ecs.ListServicesInput{
			Cluster: &cluster,
		},
	)
	if err != nil {
		return nil, err
	}

	describeServicesOutput, err := client.DescribeServices(context.TODO(),
		&ecs.DescribeServicesInput{
			Cluster:  &cluster,
			Services: listServicesOutput.ServiceArns,
		},
	)

	return describeServicesOutput.Services, err
}
