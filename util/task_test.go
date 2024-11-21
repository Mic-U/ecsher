package util

import (
	"testing"

	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func TestFilterTasksByNames(t *testing.T) {
	cases := []struct {
		tasks []string
		names []string
	}{
		{
			tasks: []string{
				"arn:aws:ecs:region:account:task/cluster-name/abcdefg",
				"arn:aws:ecs:region:account:task/cluster-name/1234567",
			},
			names: []string{
				"abcdefg",
			},
		},
		{
			tasks: []string{
				"arn:aws:ecs:region:account:task/cluster-name/abcdefg",
				"arn:aws:ecs:region:account:task/cluster-name/1234567",
			},
			names: []string{
				"1234567",
			},
		},
		{
			tasks: []string{
				"arn:aws:ecs:region:account:task/cluster-name/abcdefg",
				"arn:aws:ecs:region:account:task/cluster-name/1234567",
			},
			names: []string{},
		},
	}

	for _, c := range cases {
		result := FilterTasksByNames(c.tasks, c.names)
		if len(c.names) == 0 {
			if len(result) != 2 {
				t.Fatal("len(result) should be 2")
			}
			continue
		}
		if len(result) != len(c.names) {
			t.Fatal("len(result) should be 1")
		}
		resultName := ArnToName(result[0])
		if resultName != c.names[0] {
			t.Errorf("expect %v, got %v", c.names[0], resultName)
		}
	}
}

func TestGetCapacityProviderName(t *testing.T) {
	testTaskName := "test"
	testCapacityProvider := "FARGATE"
	cases := []struct {
		Task     ecsTypes.Task
		Expected string
	}{
		{
			Task: ecsTypes.Task{
				TaskArn: &testTaskName,
			},
			Expected: "",
		},
		{
			Task: ecsTypes.Task{
				TaskArn:              &testTaskName,
				CapacityProviderName: &testCapacityProvider,
			},
			Expected: "FARGATE",
		},
	}
	for _, c := range cases {
		actual := GetCapacityProviderName(c.Task)
		if actual != c.Expected {
			t.Fatalf("Expected is %s, but actual is %s", c.Expected, actual)
		}
	}
}
