package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func GetCluster(names []string) ([]ecsTypes.Cluster, error) {
	client := GetClient()
	if len(names) == 0 {
		listClustersOutput, err := client.ListClusters(context.TODO(), &ecs.ListClustersInput{})
		if err != nil {
			return nil, err
		}
		return DescribeCluster(listClustersOutput.ClusterArns)
	} else {
		return DescribeCluster(names)
	}
}

func DescribeCluster(names []string) ([]ecsTypes.Cluster, error) {
	client := GetClient()
	describeClustersOutput, err := client.DescribeClusters(
		context.TODO(),
		&ecs.DescribeClustersInput{
			Clusters: names,
		},
	)
	return describeClustersOutput.Clusters, err
}
