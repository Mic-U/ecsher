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
		clusterArns := []string{}
		paginater := ecs.NewListClustersPaginator(client, &ecs.ListClustersInput{})
		for paginater.HasMorePages() {
			output, err := paginater.NextPage(context.TODO())
			if err != nil {
				return []ecsTypes.Cluster{}, nil
			}
			clusterArns = append(clusterArns, output.ClusterArns...)
		}
		if len(clusterArns) == 0 {
			return []ecsTypes.Cluster{}, nil
		}
		return DescribeCluster(region, clusterArns)
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
	if err != nil {
		return []ecsTypes.Cluster{}, err
	}
	return describeClustersOutput.Clusters, err
}
