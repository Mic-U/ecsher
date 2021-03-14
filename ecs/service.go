package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// GetService returns Service list. If names is not specified, calls listServices API
func GetService(region string, cluster string, names []string) ([]ecsTypes.Service, error) {
	client := GetClient(region)
	if len(names) == 0 {
		listServicesOutput, err := client.ListServices(context.TODO(),
			&ecs.ListServicesInput{
				Cluster: &cluster,
			},
		)
		if err != nil {
			return nil, err
		}
		if len(listServicesOutput.ServiceArns) == 0 {
			return []ecsTypes.Service{}, nil
		}
		return DescribeService(region, cluster, listServicesOutput.ServiceArns)
	}
	return DescribeService(region, cluster, names)
}

// DescribeService returns Service list. This requires specifying service name
func DescribeService(region string, cluster string, names []string) ([]ecsTypes.Service, error) {
	client := GetClient(region)
	describeServicesOutput, err := client.DescribeServices(context.TODO(),
		&ecs.DescribeServicesInput{
			Cluster:  &cluster,
			Services: names,
		},
	)
	if err != nil {
		return []ecsTypes.Service{}, err
	}
	return describeServicesOutput.Services, err
}
