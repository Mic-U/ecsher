package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func ListFamily(region string, prefix string, status string) ([]string, error) {
	client := GetClient(region)
	input := ecs.ListTaskDefinitionFamiliesInput{}
	if prefix != "" {
		input.FamilyPrefix = &prefix
	}

	if status == "ACTIVE" || status == "INACTIVE" {
		input.Status = types.TaskDefinitionFamilyStatus(status)
	} else {
		input.Status = "ACTIVE"
	}

	result, err := client.ListTaskDefinitionFamilies(context.TODO(), &input)
	if err != nil {
		return []string{}, err
	}
	return result.Families, nil
}

func GetRevisions(region string, family string, status string) ([]string, error) {
	client := GetClient(region)
	input := ecs.ListTaskDefinitionsInput{}

	if family != "" {
		input.FamilyPrefix = &family
	}
	if status == "ACTIVE" || status == "INACTIVE" {
		input.Status = types.TaskDefinitionStatus(status)
	} else {
		input.Status = "ACTIVE"
	}

	result, err := client.ListTaskDefinitions(context.TODO(), &input)
	if err != nil {
		return []string{}, err
	}
	return result.TaskDefinitionArns, err
}

func DescribeDefinition(region string, name string) (types.TaskDefinition, error) {
	client := GetClient(region)
	definition, err := client.DescribeTaskDefinition(context.TODO(), &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &name,
	})
	if err != nil {
		return types.TaskDefinition{}, err
	}
	return *definition.TaskDefinition, err
}
