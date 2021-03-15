package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func GetInstance(region string, cluster string, names []string) ([]types.ContainerInstance, error) {
	client := GetClient(region)
	if len(names) == 0 {
		instances := []string{}
		paginater := ecs.NewListContainerInstancesPaginator(client, &ecs.ListContainerInstancesInput{
			Cluster: aws.String(cluster),
		})

		if paginater.HasMorePages() {
			output, err := paginater.NextPage(context.TODO())
			if err != nil {
				return []types.ContainerInstance{}, err
			}
			instances = append(instances, output.ContainerInstanceArns...)
		}
		if len(instances) == 0 {
			return []types.ContainerInstance{}, nil
		}
		return DescribeInstance(region, cluster, instances)
	}
	return DescribeInstance(region, cluster, names)

}

func DescribeInstance(region string, cluster string, names []string) ([]types.ContainerInstance, error) {
	client := GetClient(region)
	describeContainerInstanceOutput, err := client.DescribeContainerInstances(context.TODO(),
		&ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(cluster),
			ContainerInstances: names,
		},
	)
	if err != nil {
		return []types.ContainerInstance{}, nil
	}
	return describeContainerInstanceOutput.ContainerInstances, err
}
