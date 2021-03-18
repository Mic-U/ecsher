package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSInstanceClient interface {
	DescribeContainerInstances(context.Context, *ecs.DescribeContainerInstancesInput, ...func(*ecs.Options)) (*ecs.DescribeContainerInstancesOutput, error)
	ListContainerInstances(context.Context, *ecs.ListContainerInstancesInput, ...func(*ecs.Options)) (*ecs.ListContainerInstancesOutput, error)
}

type ListInstancesPager interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*ecs.Options)) (*ecs.ListContainerInstancesOutput, error)
}

func GetInstance(client ECSInstanceClient, cluster string, names []string) ([]types.ContainerInstance, error) {
	if len(names) == 0 {
		paginator := ecs.NewListContainerInstancesPaginator(client, &ecs.ListContainerInstancesInput{
			Cluster: aws.String(cluster),
		})


		instances, err := ListAllInstances(context.TODO(), paginator)
		if err != nil {
			return []types.ContainerInstance{}, err
		}
		if len(instances) == 0 {
			return []types.ContainerInstance{}, nil
		}
		return DescribeInstance(client, cluster, instances)
	}
	return DescribeInstance(client, cluster, names)

}

func ListAllInstances(ctx context.Context, paginator ListInstancesPager) ([]string, error) {
	instances := []string{}
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return []string{}, err
		}
		instances = append(instances, output.ContainerInstanceArns...)
	}
	return instances, nil
}

func DescribeInstance(client ECSInstanceClient, cluster string, names []string) ([]types.ContainerInstance, error) {
	describeContainerInstanceOutput, err := client.DescribeContainerInstances(context.TODO(),
		&ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(cluster),
			ContainerInstances: names,
		},
	)
	if err != nil {
		return []types.ContainerInstance{}, err
	}
	return describeContainerInstanceOutput.ContainerInstances, err
}
