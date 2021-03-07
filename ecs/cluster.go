package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// GetCluster returns Cluster list. If names is not specified, calls listClusters API
func GetCluster(region string, names []string) ([]ecsTypes.Cluster, error) {
	client := GetClient(region)
	if len(names) == 0 {
		listClustersOutput, err := client.ListClusters(context.TODO(), &ecs.ListClustersInput{})
		if err != nil {
			return nil, err
		}
		return DescribeCluster(region, listClustersOutput.ClusterArns)
	}
	return DescribeCluster(region, names)
}

// DescribeCluster returns Cluster list. This requires specifying cluster name
func DescribeCluster(region string, names []string) ([]ecsTypes.Cluster, error) {
	client := GetClient(region)
	describeClustersOutput, err := client.DescribeClusters(
		context.TODO(),
		&ecs.DescribeClustersInput{
			Clusters: names,
		},
	)
	return describeClustersOutput.Clusters, err
}
