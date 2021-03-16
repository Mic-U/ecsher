package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSClusterClient interface {
	DescribeClusters(context.Context, *ecs.DescribeClustersInput, ...func(*ecs.Options)) (*ecs.DescribeClustersOutput, error)
	ListClusters(context.Context, *ecs.ListClustersInput, ...func(*ecs.Options)) (*ecs.ListClustersOutput, error)
}

type ListClustersPager interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*ecs.Options)) (*ecs.ListClustersOutput, error)
}

// GetCluster returns Cluster list. If names is not specified, calls listClusters API
func GetCluster(client ECSClusterClient, names []string) ([]ecsTypes.Cluster, error) {
	if len(names) == 0 {
		paginator := ecs.NewListClustersPaginator(client, &ecs.ListClustersInput{})
		clusterArns, err := ListAllClusters(context.TODO(), paginator)
		if err != nil {
			return []ecsTypes.Cluster{}, err
		}
		if len(clusterArns) == 0 {
			return []ecsTypes.Cluster{}, nil
		}
		return DescribeCluster(client, clusterArns)
	}
	return DescribeCluster(client, names)
}

func ListAllClusters(ctx context.Context, paginator ListClustersPager) ([]string, error) {
	clusterArns := []string{}
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			return []string{}, err
		}
		clusterArns = append(clusterArns, output.ClusterArns...)
	}
	return clusterArns, nil
}

// DescribeCluster returns Cluster list. This requires specifying cluster name
func DescribeCluster(client ECSClusterClient, names []string) ([]ecsTypes.Cluster, error) {
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
