package ecs

import (
	"context"

	"github.com/Mic-U/ecsher/util"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func GetTask(cluster string, service string, names []string) ([]ecsTypes.Task, error) {
	if service == "" {
		return GetTaskInCluster(cluster, names)
	} else {
		return GetTaskInService(cluster, service, names)
	}
}

func GetTaskInCluster(cluster string, names []string) ([]ecsTypes.Task, error) {
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
	taskArns := filterTaskArns(listTasksOutput.TaskArns, names)
	if len(taskArns) == 0 {
		result := []ecsTypes.Task{}
		return result, nil
	}
	return DescribeTask(cluster, taskArns)
}

func filterTaskArns(taskArns []string, names []string) []string {
	if len(names) == 0 {
		return taskArns
	}
	var filteredTaskArns []string
	for _, taskArn := range taskArns {
		for _, name := range names {
			if taskArn == name || util.ArnToName(taskArn) == name {
				filteredTaskArns = append(filteredTaskArns, taskArn)
			}
		}
	}
	return filteredTaskArns
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
