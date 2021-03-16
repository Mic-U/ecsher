package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSServiceClient interface {
	DescribeServices(context.Context, *ecs.DescribeServicesInput, ...func(*ecs.Options)) (*ecs.DescribeServicesOutput, error)
	ListServices(context.Context, *ecs.ListServicesInput, ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
}

type ListServicesPager interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*ecs.Options)) (*ecs.ListServicesOutput, error)
}

// GetService returns Service list. If names is not specified, calls listServices API
func GetService(client ECSServiceClient, cluster string, names []string) ([]ecsTypes.Service, error) {
	if len(names) == 0 {
		paginator := ecs.NewListServicesPaginator(client, &ecs.ListServicesInput{
			Cluster: aws.String(cluster),
		})
		serviceArns, err := ListAllServices(context.TODO(), paginator)
		if err != nil {
			return []ecsTypes.Service{}, err
		}
		if len(serviceArns) == 0 {
			return []ecsTypes.Service{}, nil
		}
		return DescribeService(client, cluster, serviceArns)
	}
	return DescribeService(client, cluster, names)
}

func ListAllServices(ctx context.Context, paginator ListServicesPager) ([]string, error) {
	serviceArns := []string{}
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return []string{}, err
		}
		serviceArns = append(serviceArns, output.ServiceArns...)
	}
	return serviceArns, nil
}

// DescribeService returns Service list. This requires specifying service name
func DescribeService(client ECSServiceClient, cluster string, names []string) ([]ecsTypes.Service, error) {
	describeServicesOutput, err := client.DescribeServices(context.TODO(),
		&ecs.DescribeServicesInput{
			Cluster:  aws.String(cluster),
			Services: names,
		},
	)
	if err != nil {
		return []ecsTypes.Service{}, err
	}
	return describeServicesOutput.Services, err
}
