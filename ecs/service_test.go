package ecs

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockServiceClient struct {
	listServicesOutput ecs.ListServicesOutput
}

func (m mockServiceClient) DescribeServices(ctx context.Context, params *ecs.DescribeServicesInput, optFns ...func(*ecs.Options)) (*ecs.DescribeServicesOutput, error) {
	names := params.Services
	if len(names) == 0 || len(*params.Cluster) == 0 {
		return nil, errors.New("required param is missing")
	}
	services := []ecsTypes.Service{}
	for _, s := range names {
		service := &ecsTypes.Service{
			ServiceName: &s,
			ClusterArn:  params.Cluster,
		}
		services = append(services, *service)
	}
	result := &ecs.DescribeServicesOutput{
		Services: services,
	}
	return result, nil
}

func (m mockServiceClient) ListServices(ctx context.Context, params *ecs.ListServicesInput, optFns ...func(*ecs.Options)) (*ecs.ListServicesOutput, error) {
	return &m.listServicesOutput, nil
}

type mockServicePager struct {
	PageNum int
	Pages   []*ecs.ListServicesOutput
}

func (m *mockServicePager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}

func (m *mockServicePager) NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (output *ecs.ListServicesOutput, err error) {
	if m.PageNum >= len(m.Pages) {
		return nil, errors.New("no more pages")
	}
	output = m.Pages[m.PageNum]
	m.PageNum++
	return output, nil
}

func TestDescriveService(t *testing.T) {
	testName := "test"
	cases := []struct {
		Expected []ecsTypes.Service
		cluster  string
		names    []string
	}{
		{
			Expected: []ecsTypes.Service{
				{
					ServiceName: &testName,
				},
			},
			cluster: "test",
			names:   []string{"test"},
		}, {
			Expected: []ecsTypes.Service{
				{
					ServiceName: &testName,
				},
			},
			cluster: "",
			names:   []string{"test"},
		}, {
			Expected: []ecsTypes.Service{
				{
					ServiceName: &testName,
				},
			},
			cluster: "test",
			names:   []string{""},
		},
	}

	for i, c := range cases {
		client := &mockServiceClient{}
		result, err := DescribeService(client, c.cluster, c.names)
		if len(c.cluster) == 0 || len(c.names) == 0 {
			if err == nil {
				t.Fatal("Should return error")
			}
			continue
		}
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		if a, e := len(result), len(c.Expected); a != e {
			t.Fatalf("%d, expected %d messages, got %d", i, e, a)
		}
	}
}

func TestListAllServices(t *testing.T) {
	paginator := &mockServicePager{
		PageNum: 0,
		Pages: []*ecs.ListServicesOutput{
			{
				ServiceArns: []string{"a", "b"},
			},
			{
				ServiceArns: []string{"c", "d"},
			},
			{
				ServiceArns: []string{"e", "f"},
			},
		},
	}
	serviceArns, err := ListAllServices(context.TODO(), paginator)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if expect, actual := int(6), len(serviceArns); expect != actual {
		t.Errorf("expect %v, got %v", expect, actual)
	}

	expected := []string{"a", "b", "c", "d", "e", "f"}

	for i, actual := range serviceArns {
		if actual != expected[i] {
			t.Errorf("expect %v, got %v", expected[i], actual)
		}
	}
}

func TestGetService(t *testing.T) {
	cases := []struct {
		client  mockServiceClient
		cluster string
		names   []string
	}{
		{
			client: mockServiceClient{
				listServicesOutput: ecs.ListServicesOutput{
					NextToken:   nil,
					ServiceArns: []string{"test1", "test2", "test3"},
				},
			},
			cluster: "test",
			names:   []string{"test1", "test2"},
		},
		{
			client: mockServiceClient{
				listServicesOutput: ecs.ListServicesOutput{
					NextToken:   nil,
					ServiceArns: []string{"test1", "test2", "test3"},
				},
			},
			cluster: "test",
			names:   []string{},
		},
		{
			client: mockServiceClient{
				listServicesOutput: ecs.ListServicesOutput{
					NextToken:   nil,
					ServiceArns: []string{},
				},
			},
			cluster: "test",
			names:   []string{},
		},
	}

	for _, c := range cases {
		result, err := GetService(c.client, c.cluster, c.names)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
		if len(c.names) > 0 && len(c.client.listServicesOutput.ServiceArns) > 0 {
			if expected, actual := len(c.names), len(result); expected != actual {
				t.Errorf("expect %v, got %v", expected, actual)
			}
		} else if len(c.client.listServicesOutput.ServiceArns) > 0 {
			// When names not specified, mockListCLuster should return 3 clusters
			if expected, actual := 3, len(result); expected != actual {
				t.Errorf("expect %v, got %v", expected, actual)
			}
		} else if len(c.client.listServicesOutput.ServiceArns) == 0 {
			if expected, actual := 0, len(result); expected != actual {
				t.Errorf("expect %v, got %v", expected, actual)
			}
		}
	}
}
