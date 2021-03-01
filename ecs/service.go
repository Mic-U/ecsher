package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func GetService(cluster string, names []string) ([]ecsTypes.Service, error) {
	client := GetClient()
	if len(names) == 0 {
		listServicesOutput, err := client.ListServices(context.TODO(),
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
	client := GetClient()
	describeServicesOutput, err := client.DescribeServices(context.TODO(),
		&ecs.DescribeServicesInput{
			Cluster:  &cluster,
			Services: names,
		},
	)
	return describeServicesOutput.Services, err
}
