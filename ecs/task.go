package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func GetTask(cluster string, names []string) ([]ecsTypes.Task, error) {
	client := GetClient()
	if len(names) == 0 {
		listTasksOutput, err := client.ListTasks(context.TODO(),
			&ecs.ListTasksInput{
				Cluster: &cluster,
			},
		)
		if err != nil {
			return nil, err
		}
		return DescribeTask(cluster, listTasksOutput.TaskArns)
	} else {
		return DescribeTask(cluster, names)
	}
}

func GetTaskInService(cluster string, service string, names []string) ([]ecsTypes.Task, error) {
	client := GetClient()
	listTasksOutput, err := client.ListTasks(context.TODO(),
		&ecs.ListTasksInput{
			Cluster:     &cluster,
			ServiceName: &service,
		},
	)
	if err != nil {
		return nil, err
	}
	return DescribeTask(cluster, listTasksOutput.TaskArns)
}

func DescribeTask(cluster string, names []string) ([]ecsTypes.Task, error) {
	client := GetClient()
	describeTasksOutput, err := client.DescribeTasks(context.TODO(),
		&ecs.DescribeTasksInput{
			Cluster: &cluster,
			Tasks:   names,
		},
	)
	return describeTasksOutput.Tasks, err
}
