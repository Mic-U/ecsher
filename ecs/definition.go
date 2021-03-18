package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSTaskDefinitionClient interface {
	DescribeTaskDefinition(context.Context, *ecs.DescribeTaskDefinitionInput, ...func(*ecs.Options)) (*ecs.DescribeTaskDefinitionOutput, error)
	ListTaskDefinitionFamilies(context.Context, *ecs.ListTaskDefinitionFamiliesInput, ...func(*ecs.Options)) (*ecs.ListTaskDefinitionFamiliesOutput, error)
	ListTaskDefinitions(context.Context, *ecs.ListTaskDefinitionsInput, ...func(*ecs.Options)) (*ecs.ListTaskDefinitionsOutput, error)
}

type ListTaskDefinitionFamiliesPager interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*ecs.Options)) (*ecs.ListTaskDefinitionFamiliesOutput, error)
}

type ListTaskDefinitionsPager interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*ecs.Options)) (*ecs.ListTaskDefinitionsOutput, error)
}

func GetFamily(client ECSTaskDefinitionClient, prefix string, status string) ([]string, error) {
	input := ecs.ListTaskDefinitionFamiliesInput{}
	if prefix != "" {
		input.FamilyPrefix = aws.String(prefix)
	}

	if status == "ACTIVE" || status == "INACTIVE" {
		input.Status = types.TaskDefinitionFamilyStatus(status)
	} else {
		input.Status = "ACTIVE"
	}

	paginator := ecs.NewListTaskDefinitionFamiliesPaginator(client, &input)
	return ListAllFamilies(context.TODO(), paginator)
}

func ListAllFamilies(ctx context.Context, paginator ListTaskDefinitionFamiliesPager) ([]string, error) {
	families := []string{}
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return families, err
		}
		families = append(families, output.Families...)
	}
	return families, nil
}

func GetRevision(client ECSTaskDefinitionClient, family string, status string) ([]string, error) {
	input := ecs.ListTaskDefinitionsInput{}

	if family != "" {
		input.FamilyPrefix = aws.String(family)
	}
	if status == "ACTIVE" || status == "INACTIVE" {
		input.Status = types.TaskDefinitionStatus(status)
	} else {
		input.Status = "ACTIVE"
	}

	paginator := ecs.NewListTaskDefinitionsPaginator(client, &input)
	return ListAllRevisions(context.TODO(), paginator)
}

func ListAllRevisions(ctx context.Context, paginator ListTaskDefinitionsPager) ([]string, error) {
	definitions := []string{}
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return definitions, err
		}
		definitions = append(definitions, output.TaskDefinitionArns...)
	}
	return definitions, nil
}

func DescribeDefinition(client ECSTaskDefinitionClient, name string) (types.TaskDefinition, error) {
	definition, err := client.DescribeTaskDefinition(context.TODO(), &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(name),
	})
	if err != nil {
		return types.TaskDefinition{}, err
	}
	return *definition.TaskDefinition, err
}
