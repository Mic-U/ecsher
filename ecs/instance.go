package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func GetInstance(region string, cluster string, names []string) ([]types.ContainerInstance, error) {
	client := GetClient(region)
	if len(names) == 0 {
		listContainerInstancesOutput, err := client.ListContainerInstances(context.TODO(),
			&ecs.ListContainerInstancesInput{
				Cluster: &cluster,
			},
		)
		if err != nil {
			return nil, err
		}
		if len(listContainerInstancesOutput.ContainerInstanceArns) == 0 {
			return []types.ContainerInstance{}, nil
		}
		return DescribeInstance(region, cluster, listContainerInstancesOutput.ContainerInstanceArns)
	}
	return DescribeInstance(region, cluster, names)

}

func DescribeInstance(region string, cluster string, names []string) ([]types.ContainerInstance, error) {
	client := GetClient(region)
	describeContainerInstanceOutput, err := client.DescribeContainerInstances(context.TODO(),
		&ecs.DescribeContainerInstancesInput{
			Cluster:            &cluster,
			ContainerInstances: names,
		},
	)
	if err != nil {
		return []types.ContainerInstance{}, nil
	}
	return describeContainerInstanceOutput.ContainerInstances, err
}
