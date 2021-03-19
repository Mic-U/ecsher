package util

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func TestListContainerNames(t *testing.T) {
	testContainers := []ecsTypes.ContainerDefinition{
		{
			Name: aws.String("a"),
		},
		{
			Name: aws.String("b"),
		},
		{
			Name: aws.String("c"),
		},
	}
	result := ListContainerNames(testContainers)
	for i, r := range result {
		if r != *testContainers[i].Name {
			t.Errorf("expected is %s, but actual is %s", *testContainers[i].Name, r)
		}
	}
}

func TestChooseContainer(t *testing.T) {
	testContainers := []ecsTypes.ContainerDefinition{
		{
			Name: aws.String("a"),
		},
		{
			Name: aws.String("b"),
		},
		{
			Name: aws.String("c"),
		},
	}
	result1, err := ChooseContainer(testContainers, "b")
	if err != nil {
		t.Error(err)
	}
	if *result1.Name != "b" {
		t.Errorf("expected is b, but actual is %s", *result1.Name)
	}

	_, err = ChooseContainer(testContainers, "d")
	if err == nil {
		t.Fatal("Should return error")
	}
}
