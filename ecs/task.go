package ecs

import (
	"context"

	"github.com/Mic-U/ecsher/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ECSTaskClient interface {
	DescribeTasks(context.Context, *ecs.DescribeTasksInput, ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error)
	ExecuteCommand(context.Context, *ecs.ExecuteCommandInput, ...func(*ecs.Options)) (*ecs.ExecuteCommandOutput, error)
	ListTasks(context.Context, *ecs.ListTasksInput, ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
}

type ListTasksPager interface {
	HasMorePages() bool
	NextPage(context.Context, ...func(*ecs.Options)) (*ecs.ListTasksOutput, error)
}

// GetTask returns Cluster list. If Service is not specified, returns tasks in the service
func GetTask(client ECSTaskClient, cluster string, service string, names []string) ([]ecsTypes.Task, error) {
	if service == "" {
		return GetTaskInCluster(client, cluster, names)
	}
	return GetTaskInService(client, cluster, service, names)
}

// GetTaskInCluster returns tasks in the Cluster
func GetTaskInCluster(client ECSTaskClient, cluster string, names []string) ([]ecsTypes.Task, error) {
	if len(names) == 0 {
		paginator := ecs.NewListTasksPaginator(client, &ecs.ListTasksInput{
			Cluster: aws.String(cluster),
		})
		tasks, err := ListAllTasks(context.TODO(), paginator)
		if err != nil {
			return []ecsTypes.Task{}, err
		}
		if len(tasks) == 0 {
			return []ecsTypes.Task{}, nil
		}
		return DescribeTask(client, cluster, tasks)
	}
	return DescribeTask(client, cluster, names)
}

// GetTaskInService returns tasks in the Service in the Cluster
func GetTaskInService(client ECSTaskClient, cluster string, service string, names []string) ([]ecsTypes.Task, error) {
	paginator := ecs.NewListTasksPaginator(client, &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	})
	tasks, err := ListAllTasks(context.TODO(), paginator)
	if err != nil {
		return []ecsTypes.Task{}, err
	}
	filtered := util.FilterTasksByNames(tasks, names)
	if len(filtered) == 0 {
		result := []ecsTypes.Task{}
		return result, nil
	}
	return DescribeTask(client, cluster, filtered)
}

func ListAllTasks(ctx context.Context, paginator ListTasksPager) ([]string, error) {
	tasks := []string{}
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return []string{}, err
		}
		tasks = append(tasks, output.TaskArns...)
	}
	return tasks, nil
}

// DescribeTask returns Task list. This requires specifying task name
func DescribeTask(client ECSTaskClient, cluster string, names []string) ([]ecsTypes.Task, error) {
	describeTasksOutput, err := client.DescribeTasks(context.TODO(),
		&ecs.DescribeTasksInput{
			Cluster: aws.String(cluster),
			Tasks:   names,
		},
	)
	if err != nil {
		return []ecsTypes.Task{}, err
	}
	return describeTasksOutput.Tasks, err
}

type ExecuteCmmandParams struct {
	Task        string
	Container   string
	Command     string
	Interactive bool
}

func ExecuteCommand(client ECSTaskClient, cluster string, params *ExecuteCmmandParams) (*ecs.ExecuteCommandOutput, error) {
	input := &ecs.ExecuteCommandInput{
		Cluster:     aws.String(cluster),
		Task:        aws.String(params.Task),
		Interactive: params.Interactive,
		Command:     aws.String(params.Command),
	}
	if len(params.Container) > 0 {
		input.Container = aws.String(params.Container)
	}
	executeCommandOutput, err := client.ExecuteCommand(context.TODO(), input)
	return executeCommandOutput, err
}
