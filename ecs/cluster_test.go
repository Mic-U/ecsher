package ecs

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockClusterClient struct {
	listClustersOutput ecs.ListClustersOutput
}

func (m mockClusterClient) DescribeClusters(ctx context.Context, params *ecs.DescribeClustersInput, optFns ...func(*ecs.Options)) (*ecs.DescribeClustersOutput, error) {
	names := params.Clusters
	clusters := []ecsTypes.Cluster{}
	for _, c := range names {
		cluster := &ecsTypes.Cluster{
			ClusterName: &c,
		}
		clusters = append(clusters, *cluster)
	}
	result := &ecs.DescribeClustersOutput{
		Clusters: clusters,
	}
	return result, nil
}

func (m mockClusterClient) ListClusters(context.Context, *ecs.ListClustersInput, ...func(*ecs.Options)) (*ecs.ListClustersOutput, error) {
	return &m.listClustersOutput, nil
}

type mockClusterPager struct {
	PageNum int
	Pages   []*ecs.ListClustersOutput
}

func (m *mockClusterPager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}

func (m *mockClusterPager) NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (output *ecs.ListClustersOutput, err error) {
	if m.PageNum >= len(m.Pages) {
		return nil, fmt.Errorf("no more pages")
	}
	output = m.Pages[m.PageNum]
	m.PageNum++
	fmt.Println(m.PageNum)
	return output, nil
}

func TestDescribeCluster(t *testing.T) {
	testName := "test"
	cases := []struct {
		Expected []ecsTypes.Cluster
		names    []string
	}{
		{
			Expected: []ecsTypes.Cluster{
				{
					ClusterName: &testName,
				},
			},
			names: []string{testName},
		},
	}
	for i, c := range cases {
		client := mockClusterClient{}
		result, err := DescribeCluster(client, c.names)

		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}

		if a, e := len(result), len(c.Expected); a != e {
			t.Fatalf("%d, expected %d messages, got %d", i, e, a)
		}
	}
}

func TestListAllClusters(t *testing.T) {
	paginator := &mockClusterPager{
		PageNum: 0,
		Pages: []*ecs.ListClustersOutput{
			{
				ClusterArns: []string{"a", "b"},
			},
			{
				ClusterArns: []string{"c", "d"},
			},
			{
				ClusterArns: []string{"e", "f"},
			},
		},
	}
	clusterArns, err := ListAllClusters(context.TODO(), paginator)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if expect, actual := int(6), len(clusterArns); expect != actual {
		t.Errorf("expect %v, got %v", expect, actual)
	}
}

func TestGetCluster(t *testing.T) {

	cases := []struct {
		client mockClusterClient
		names  []string
	}{
		{
			client: mockClusterClient{
				listClustersOutput: ecs.ListClustersOutput{
					ClusterArns: []string{"test1", "test2", "test3"},
					NextToken:   nil,
				},
			},
			names: []string{"test1", "test2"},
		},
		{
			client: mockClusterClient{
				listClustersOutput: ecs.ListClustersOutput{
					ClusterArns: []string{"test1", "test2", "test3"},
					NextToken:   nil,
				},
			},
			names: []string{},
		},
		{
			client: mockClusterClient{
				listClustersOutput: ecs.ListClustersOutput{
					ClusterArns: []string{},
					NextToken:   nil,
				},
			},
			names: []string{},
		},
	}

	for _, c := range cases {
		result, err := GetCluster(c.client, c.names)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
		if len(c.names) > 0 && len(c.client.listClustersOutput.ClusterArns) > 0 {
			if expected, actual := len(c.names), len(result); expected != actual {
				t.Errorf("expect %v, got %v", expected, actual)
			}
		} else if len(c.client.listClustersOutput.ClusterArns) > 0 {
			// When names not specified, mockListCLuster should return 3 clusters
			if expected, actual := 3, len(result); expected != actual {
				t.Errorf("expect %v, got %v", expected, actual)
			}
		} else if len(c.client.listClustersOutput.ClusterArns) == 0 {
			if expected, actual := 0, len(result); expected != actual {
				t.Errorf("expect %v, got %v", expected, actual)
			}
		}
	}
}
