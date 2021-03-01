package ecs

import (
	"context"

	util "github.com/Mic-U/kecs/util"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func GetCluster() ([]ecsTypes.Cluster, error) {
	client := GetClient()
	listClustersOutput, err := client.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		return nil, err
	}

	clusterNames := make([]string, len(listClustersOutput.ClusterArns))
	for i, clusterArn := range listClustersOutput.ClusterArns {
		clusterNames[i] = util.ARNtoName(clusterArn)
	}
	describeClustersOutput, err := client.DescribeClusters(
		context.TODO(),
		&ecs.DescribeClustersInput{
			Clusters: clusterNames,
		},
	)
	return describeClustersOutput.Clusters, nil
}
