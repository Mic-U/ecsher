package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func ListFamily(region string, prefix string, status string) ([]string, error) {
	client := GetClient(region)
	input := ecs.ListTaskDefinitionFamiliesInput{}
	if prefix != "" {
		input.FamilyPrefix = aws.String(prefix)
	}

	if status == "ACTIVE" || status == "INACTIVE" {
		input.Status = types.TaskDefinitionFamilyStatus(status)
	} else {
		input.Status = "ACTIVE"
	}

	families := []string{}
	paginater := ecs.NewListTaskDefinitionFamiliesPaginator(client, &input)
	for paginater.HasMorePages() {
		output, err := paginater.NextPage(context.TODO())
		if err != nil {
			return families, err
		}
		families = append(families, output.Families...)
	}
	return families, nil
}

func GetRevisions(region string, family string, status string) ([]string, error) {
	client := GetClient(region)
	input := ecs.ListTaskDefinitionsInput{}

	if family != "" {
		input.FamilyPrefix = aws.String(family)
	}
	if status == "ACTIVE" || status == "INACTIVE" {
		input.Status = types.TaskDefinitionStatus(status)
	} else {
		input.Status = "ACTIVE"
	}

	definitions := []string{}
	paginater := ecs.NewListTaskDefinitionsPaginator(client, &input)
	for paginater.HasMorePages() {
		output, err := paginater.NextPage(context.TODO())
		if err != nil {
			return definitions, err
		}
		definitions = append(definitions, output.TaskDefinitionArns...)
	}
	return definitions, nil
}

func DescribeDefinition(region string, name string) (types.TaskDefinition, error) {
	client := GetClient(region)
	definition, err := client.DescribeTaskDefinition(context.TODO(), &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(name),
	})
	if err != nil {
		return types.TaskDefinition{}, err
	}
	return *definition.TaskDefinition, err
}
