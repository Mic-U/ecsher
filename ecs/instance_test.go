package ecs

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockInstanceClient struct {
	listInstancesOutput ecs.ListContainerInstancesOutput
}

func (m mockInstanceClient) DescribeContainerInstances(ctx context.Context, params *ecs.DescribeContainerInstancesInput, optFns ...func(*ecs.Options)) (*ecs.DescribeContainerInstancesOutput, error) {
	names := params.ContainerInstances
	if len(names) == 0 || len(*params.Cluster) == 0 {
		return nil, errors.New("required param is missing")
	}
	instances := []ecsTypes.ContainerInstance{}
	for _, s := range names {
		instance := &ecsTypes.ContainerInstance{
			ContainerInstanceArn: &s,
		}
		instances = append(instances, *instance)
	}
	result := &ecs.DescribeContainerInstancesOutput{
		ContainerInstances: instances,
	}
	return result, nil
}

func (m mockInstanceClient) ListContainerInstances(ctx context.Context, params *ecs.ListContainerInstancesInput, optFns ...func(*ecs.Options)) (*ecs.ListContainerInstancesOutput, error) {
	if len(*params.Cluster) == 0 {
		return nil, errors.New("required param is missing")
	}
	return &m.listInstancesOutput, nil
}

type mockInstancePager struct {
	PageNum int
	Pages   []*ecs.ListContainerInstancesOutput
}

func (m *mockInstancePager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}

func (m *mockInstancePager) NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (output *ecs.ListContainerInstancesOutput, err error) {
	if m.PageNum >= len(m.Pages) {
		return nil, errors.New("no more pages")
	}
	output = m.Pages[m.PageNum]
	m.PageNum++
	return output, nil
}

func TestDescribeInstance(t *testing.T) {
	testName := "test"
	cases := []struct {
		Expected []ecsTypes.ContainerInstance
		cluster  string
		names    []string
	}{
		{
			Expected: []ecsTypes.ContainerInstance{
				{
					ContainerInstanceArn: &testName,
				},
			},
			cluster: "test",
			names:   []string{"test"},
		},
		{
			Expected: []ecsTypes.ContainerInstance{
				{
					ContainerInstanceArn: &testName,
				},
			},
			cluster: "",
			names:   []string{"test"},
		},
		{
			Expected: []ecsTypes.ContainerInstance{
				{
					ContainerInstanceArn: &testName,
				},
			},
			cluster: "test",
			names:   []string{},
		},
	}

	for i, c := range cases {
		client := &mockInstanceClient{}
		result, err := DescribeInstance(client, c.cluster, c.names)
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

func TestListAllInstances(t *testing.T) {
	paginator := &mockInstancePager{
		PageNum: 0,
		Pages: []*ecs.ListContainerInstancesOutput{
			{
				ContainerInstanceArns: []string{"a", "b"},
			},
			{
				ContainerInstanceArns: []string{"c", "d"},
			},
			{
				ContainerInstanceArns: []string{"e", "f"},
			},
		},
	}
	instanceArns, err := ListAllInstances(context.TODO(), paginator)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if expect, actual := int(6), len(instanceArns); expect != actual {
		t.Errorf("expect %v, got %v", expect, actual)
	}

	expected := []string{"a", "b", "c", "d", "e", "f"}

	for i, actual := range instanceArns {
		if actual != expected[i] {
			t.Errorf("expect %v, got %v", expected[i], actual)
		}
	}
}

func TestGetInstance(t *testing.T) {
	cases := []struct {
		client         mockInstanceClient
		cluster        string
		names          []string
		expectedLength int
	}{
		{
			client: mockInstanceClient{
				listInstancesOutput: ecs.ListContainerInstancesOutput{
					NextToken:             nil,
					ContainerInstanceArns: []string{"test", "test2"},
				},
			},
			cluster:        "test",
			names:          []string{"test"},
			expectedLength: 1,
		},
		{
			client: mockInstanceClient{
				listInstancesOutput: ecs.ListContainerInstancesOutput{
					NextToken:             nil,
					ContainerInstanceArns: []string{"test", "test2"},
				},
			},
			cluster:        "test",
			names:          []string{},
			expectedLength: 2,
		},
		{
			client: mockInstanceClient{
				listInstancesOutput: ecs.ListContainerInstancesOutput{
					NextToken:             nil,
					ContainerInstanceArns: []string{"test", "test2"},
				},
			},
			cluster:        "",
			names:          []string{"test"},
			expectedLength: 0,
		},
		{
			client: mockInstanceClient{
				listInstancesOutput: ecs.ListContainerInstancesOutput{
					NextToken:             nil,
					ContainerInstanceArns: []string{"test", "test2"},
				},
			},
			cluster:        "",
			names:          []string{},
			expectedLength: 0,
		},
		{
			client: mockInstanceClient{
				listInstancesOutput: ecs.ListContainerInstancesOutput{
					NextToken:             nil,
					ContainerInstanceArns: []string{},
				},
			},
			cluster:        "test",
			names:          []string{},
			expectedLength: 0,
		},
	}
	for _, c := range cases {
		actual, err := GetInstance(c.client, c.cluster, c.names)
		if len(c.cluster) == 0 {
			if err == nil {
				t.Fatal("Should return error")
			}
			continue
		}
		if err != nil {
			t.Fatalf("%d, unexpected error", err)
		}
		if len(actual) != c.expectedLength {
			t.Errorf("expect %v, got %v", c.expectedLength, len(actual))
		}
	}
}
