package ecs

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockDefinitionClient struct {
	listTaskDefinitionsOutput        ecs.ListTaskDefinitionsOutput
	listTaskDefinitionFamiliesOutput ecs.ListTaskDefinitionFamiliesOutput
}

func (m mockDefinitionClient) DescribeTaskDefinition(ctx context.Context, params *ecs.DescribeTaskDefinitionInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTaskDefinitionOutput, error) {
	if len(*params.TaskDefinition) == 0 {
		return nil, errors.New("required param is missing")
	}
	result := &ecs.DescribeTaskDefinitionOutput{
		TaskDefinition: &ecsTypes.TaskDefinition{
			TaskDefinitionArn: params.TaskDefinition,
		},
	}
	return result, nil
}

func (m mockDefinitionClient) ListTaskDefinitions(ctx context.Context, params *ecs.ListTaskDefinitionsInput, optFns ...func(*ecs.Options)) (*ecs.ListTaskDefinitionsOutput, error) {
	if params.Status != "ACTIVE" && params.Status != "INACTIVE" {
		return nil, errors.New("Status should be active or inactive")
	}
	return &m.listTaskDefinitionsOutput, nil
}

func (m mockDefinitionClient) ListTaskDefinitionFamilies(ctx context.Context, params *ecs.ListTaskDefinitionFamiliesInput, optFns ...func(*ecs.Options)) (*ecs.ListTaskDefinitionFamiliesOutput, error) {
	if params.Status != "ACTIVE" && params.Status != "INACTIVE" {
		return nil, errors.New("Status should be active or inactive")
	}
	return &m.listTaskDefinitionFamiliesOutput, nil
}

type mockTaskDefinitionsPager struct {
	PageNum int
	Pages   []*ecs.ListTaskDefinitionsOutput
}

func (m *mockTaskDefinitionsPager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}
func (m *mockTaskDefinitionsPager) NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (output *ecs.ListTaskDefinitionsOutput, err error) {
	if m.PageNum >= len(m.Pages) {
		return nil, fmt.Errorf("no more pages")
	}
	output = m.Pages[m.PageNum]
	m.PageNum++
	return output, nil
}

type mockTaskDefinitionFamiliesPager struct {
	PageNum int
	Pages   []*ecs.ListTaskDefinitionFamiliesOutput
}

func (m *mockTaskDefinitionFamiliesPager) HasMorePages() bool {
	return m.PageNum < len(m.Pages)
}
func (m *mockTaskDefinitionFamiliesPager) NextPage(ctx context.Context, optFns ...func(*ecs.Options)) (output *ecs.ListTaskDefinitionFamiliesOutput, err error) {
	if m.PageNum >= len(m.Pages) {
		return nil, fmt.Errorf("no more pages")
	}
	output = m.Pages[m.PageNum]
	m.PageNum++
	return output, nil
}

func TestDescribeDefinition(t *testing.T) {
	cases := []string{"", "test"}
	for _, c := range cases {
		client := &mockDefinitionClient{}
		result, err := DescribeDefinition(client, c)
		if len(c) == 0 {
			if err == nil {
				t.Fatal("Should return error")
			}
			continue
		}
		actual := *result.TaskDefinitionArn
		if actual != c {
			t.Errorf("expect %v, got %v", c, actual)
		}
	}
}

func TestListAllRevisions(t *testing.T) {
	paginator := &mockTaskDefinitionsPager{
		PageNum: 0,
		Pages: []*ecs.ListTaskDefinitionsOutput{
			{
				TaskDefinitionArns: []string{"a", "b"},
			},
			{
				TaskDefinitionArns: []string{"c", "d"},
			},
			{
				TaskDefinitionArns: []string{"e", "f"},
			},
		},
	}
	taskDefinitionArns, err := ListAllRevisions(context.TODO(), paginator)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if expect, actual := int(6), len(taskDefinitionArns); expect != actual {
		t.Errorf("expect %v, got %v", expect, actual)
	}

	expected := []string{"a", "b", "c", "d", "e", "f"}

	for i, actual := range taskDefinitionArns {
		if actual != expected[i] {
			t.Errorf("expect %v, got %v", expected[i], actual)
		}
	}
}

func TestListAllFamilies(t *testing.T) {
	paginator := &mockTaskDefinitionFamiliesPager{
		PageNum: 0,
		Pages: []*ecs.ListTaskDefinitionFamiliesOutput{
			{
				Families: []string{"a", "b"},
			},
			{
				Families: []string{"c", "d"},
			},
			{
				Families: []string{"e", "f"},
			},
		},
	}

	families, err := ListAllFamilies(context.TODO(), paginator)
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if expect, actual := int(6), len(families); expect != actual {
		t.Errorf("expect %v, got %v", expect, actual)
	}

	expected := []string{"a", "b", "c", "d", "e", "f"}

	for i, actual := range families {
		if actual != expected[i] {
			t.Errorf("expect %v, got %v", expected[i], actual)
		}
	}
}

func TestGetRevision(t *testing.T) {
	cases := []struct {
		client mockDefinitionClient
		status string
		family string
	}{
		{
			client: mockDefinitionClient{
				listTaskDefinitionsOutput: ecs.ListTaskDefinitionsOutput{
					TaskDefinitionArns: []string{"a", "b"},
				},
			},
			status: "ACTIVE",
			family: "test",
		},
		{
			client: mockDefinitionClient{
				listTaskDefinitionsOutput: ecs.ListTaskDefinitionsOutput{
					TaskDefinitionArns: []string{"a", "b"},
				},
			},
			status: "INACTIVE",
			family: "test",
		},
		{
			client: mockDefinitionClient{
				listTaskDefinitionsOutput: ecs.ListTaskDefinitionsOutput{
					TaskDefinitionArns: []string{"a", "b"},
				},
			},
			status: "ACTIVEE",
			family: "test",
		},
		{
			client: mockDefinitionClient{
				listTaskDefinitionsOutput: ecs.ListTaskDefinitionsOutput{
					TaskDefinitionArns: []string{"a", "b"},
				},
			},
			status: "ACTIVE",
			family: "",
		},
	}

	for _, c := range cases {
		_, err := GetRevision(c.client, c.family, c.status)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
	}
}

func TestGetFamily(t *testing.T) {
	cases := []struct {
		client mockDefinitionClient
		status string
		prefix string
	}{
		{
			client: mockDefinitionClient{
				listTaskDefinitionFamiliesOutput: ecs.ListTaskDefinitionFamiliesOutput{
					Families: []string{"a", "b"},
				},
			},
			status: "ACTIVE",
			prefix: "test",
		},
		{
			client: mockDefinitionClient{
				listTaskDefinitionFamiliesOutput: ecs.ListTaskDefinitionFamiliesOutput{
					Families: []string{"a", "b"},
				},
			},
			status: "INACTIVE",
			prefix: "test",
		},
		{
			client: mockDefinitionClient{
				listTaskDefinitionFamiliesOutput: ecs.ListTaskDefinitionFamiliesOutput{
					Families: []string{"a", "b"},
				},
			},
			status: "ACTIVEE",
			prefix: "test",
		},
		{
			client: mockDefinitionClient{
				listTaskDefinitionFamiliesOutput: ecs.ListTaskDefinitionFamiliesOutput{
					Families: []string{"a", "b"},
				},
			},
			status: "ACTIVE",
			prefix: "",
		},
	}
	for _, c := range cases {
		_, err := GetFamily(c.client, c.prefix, c.status)
		if err != nil {
			t.Fatalf("expect no error, got %v", err)
		}
	}
}
