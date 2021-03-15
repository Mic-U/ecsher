package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// GetTask returns Cluster list. If Service is not specified, returns tasks in the service
func GetTask(region string, cluster string, service string, names []string) ([]ecsTypes.Task, error) {
	if service == "" {
		return GetTaskInCluster(region, cluster, names)
	}
	return GetTaskInService(region, cluster, service, names)
}

// GetTaskInCluster returns tasks in the Cluster
func GetTaskInCluster(region string, cluster string, names []string) ([]ecsTypes.Task, error) {
	client := GetClient(region)
	if len(names) == 0 {
		tasks := []string{}
		paginater := ecs.NewListTasksPaginator(client, &ecs.ListTasksInput{
			Cluster: aws.String(cluster),
		})
		if paginater.HasMorePages() {
			output, err := paginater.NextPage(context.TODO())
			if err != nil {
				return []ecsTypes.Task{}, nil
			}
			tasks = append(tasks, output.TaskArns...)
		}
		if len(tasks) == 0 {
			return []ecsTypes.Task{}, nil
		}
		return DescribeTask(region, cluster, tasks)
	}
	return DescribeTask(region, cluster, names)
}

// GetTaskInService returns tasks in the Service in the Cluster
func GetTaskInService(region string, cluster string, service string, names []string) ([]ecsTypes.Task, error) {
	client := GetClient(region)
	tasks := []string{}
	paginater := ecs.NewListTasksPaginator(client, &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	})
	if paginater.HasMorePages() {
		output, err := paginater.NextPage(context.TODO())
		if err != nil {
			return []ecsTypes.Task{}, err
		}
		tasks = append(tasks, output.TaskArns...)
	}

	if len(tasks) == 0 {
		result := []ecsTypes.Task{}
		return result, nil
	}
	return DescribeTask(region, cluster, tasks)
}

// DescribeTask returns Task list. This requires specifying task name
func DescribeTask(region string, cluster string, names []string) ([]ecsTypes.Task, error) {
	client := GetClient(region)
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
