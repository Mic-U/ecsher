package ecs

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockTaskClient struct {
	listTasksOutput ecs.ListTasksOutput
}

func (m mockTaskClient) DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error) {
	names := params.Tasks
	if len(names) == 0 || len(*params.Cluster) == 0 {
		return nil, errors.New("required param is missing")
	}
	tasks := []ecsTypes.Task{}
	for _, t := range names {
		task := &ecsTypes.Task{
			ClusterArn: params.Cluster,
			TaskArn:    &t,
		}
		tasks = append(tasks, *task)
	}
	result := &ecs.DescribeTasksOutput{
		Tasks: tasks,
	}
	return result, nil
}

func (m mockTaskClient) ListTasks(ctx context.Context, params *ecs.ListTasksInput, optFns ...func(*ecs.Options)) (*ecs.ListTasksOutput, error) {
	if len(*params.Cluster) == 0 {
		return nil, errors.New("required param is missing")
	}
	return &m.listTasksOutput, nil
}

func (m mockTaskClient) ExecuteCommand(ctx context.Context, params *ecs.ExecuteCommandInput, optFns ...func(*ecs.Options)) (*ecs.ExecuteCommandOutput, error) {
	if len(*params.Cluster) == 0 || len(*params.Task) == 0 {
		return nil, errors.New("required param is missing")
	}
	return &ecs.ExecuteCommandOutput{
		ClusterArn: params.Cluster,
	}, nil
}

type mockTaskPager struct {
	PageNum int
	Pages   []*ecs.ListTasksOutput
}

func (m *mockTaskPager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}

func (m *mockTaskPager) NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (output *ecs.ListTasksOutput, err error) {
	if m.PageNum >= len(m.Pages) {
		return nil, fmt.Errorf("no more pages")
	}
	output = m.Pages[m.PageNum]
	m.PageNum++
	return output, nil
}

func TestDescribeTask(t *testing.T) {
	testName := "test"
	cases := []struct {
		Expected []ecsTypes.Task
		cluster  string
		names    []string
	}{
		{
			Expected: []ecsTypes.Task{
				{
					TaskArn: &testName,
				},
			},
			cluster: "test",
			names:   []string{"test"},
		},
		{
			Expected: []ecsTypes.Task{
				{
					TaskArn: &testName,
				},
			},
			cluster: "",
			names:   []string{"test"},
		},
		{
			Expected: []ecsTypes.Task{
				{
					TaskArn: &testName,
				},
			},
			cluster: "test",
			names:   []string{},
		},
	}

	for i, c := range cases {
		client := &mockTaskClient{}
		result, err := DescribeTask(client, c.cluster, c.names)
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

func TestListAllTasks(t *testing.T) {
	paginator := &mockTaskPager{
		PageNum: 0,
		Pages: []*ecs.ListTasksOutput{
			{
				TaskArns: []string{"a", "b"},
			},
			{
				TaskArns: []string{"c", "d"},
			},
			{
				TaskArns: []string{"e", "f"},
			},
		},
	}

	tasks, err := ListAllTasks(context.TODO(), paginator)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if expect, actual := int(6), len(tasks); expect != actual {
		t.Errorf("expect %v, got %v", expect, actual)
	}

	expected := []string{"a", "b", "c", "d", "e", "f"}

	for i, actual := range tasks {
		if actual != expected[i] {
			t.Errorf("expect %v, got %v", expected[i], actual)
		}
	}
}

func TestGetTaskInCluster(t *testing.T) {
	cases := []struct {
		client  mockTaskClient
		cluster string
		names   []string
	}{
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{"test"},
				},
			},
			cluster: "",
			names:   []string{"test"},
		},
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{"test"},
				},
			},
			cluster: "",
			names:   []string{},
		},
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{},
				},
			},
			cluster: "test",
			names:   []string{},
		},
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{"test"},
				},
			},
			cluster: "test",
			names:   []string{},
		},
	}

	for _, c := range cases {
		result, err := GetTaskInCluster(c.client, c.cluster, c.names)
		if len(c.cluster) == 0 {
			if err == nil {
				t.Fatal("Should return error")
			}
			continue
		}
		if len(c.client.listTasksOutput.TaskArns) == 0 {
			if len(result) != 0 {
				t.Fatal("should retun 0 results")
			}
			continue
		}
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}

	}
}

func TestGetTaskinService(t *testing.T) {
	cases := []struct {
		client  mockTaskClient
		cluster string
		service string
		names   []string
	}{
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{"test"},
				},
			},
			cluster: "",
			service: "test",
			names:   []string{"test"},
		},
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{"test"},
				},
			},
			cluster: "",
			service: "test",
			names:   []string{},
		},
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{},
				},
			},
			cluster: "test",
			service: "test",
			names:   []string{},
		},
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{"test"},
				},
			},
			cluster: "test",
			service: "test",
			names:   []string{},
		},
	}

	for _, c := range cases {
		result, err := GetTaskInService(c.client, c.cluster, c.service, c.names)
		if len(c.cluster) == 0 {
			if err == nil {
				t.Fatal("Should return error")
			}
			continue
		}
		if len(c.client.listTasksOutput.TaskArns) == 0 {
			if len(result) != 0 {
				t.Fatal("should retun 0 results")
			}
			continue
		}
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
	}
}

func TestGetTask(t *testing.T) {
	cases := []struct {
		client  mockTaskClient
		cluster string
		service string
		names   []string
	}{
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{"test"},
				},
			},
			cluster: "test",
			service: "",
			names:   []string{"test"},
		},
		{
			client: mockTaskClient{
				listTasksOutput: ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []string{"test"},
				},
			},
			cluster: "test",
			service: "test",
			names:   []string{"test"},
		},
	}

	for _, c := range cases {
		_, err := GetTask(c.client, c.cluster, c.service, c.names)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
	}
}
