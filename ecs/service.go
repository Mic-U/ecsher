package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

var Client = GetClient()

func GetService(cluster string, names []string) ([]ecsTypes.Service, error) {
	if len(names) == 0 {
		listServicesOutput, err := Client.ListServices(context.TODO(),
			&ecs.ListServicesInput{
				Cluster: &cluster,
			},
		)
		if err != nil {
			return nil, err
		}
		return DescribeService(cluster, listServicesOutput.ServiceArns)
	} else {
		return DescribeService(cluster, names)
	}
}

func DescribeService(cluster string, names []string) ([]ecsTypes.Service, error) {
	describeServicesOutput, err := Client.DescribeServices(context.TODO(),
		&ecs.DescribeServicesInput{
			Cluster:  &cluster,
			Services: names,
		},
	)
	return describeServicesOutput.Services, err
}
