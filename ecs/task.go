package ecs

import (
	"context"

	"github.com/Mic-U/ecsher/util"
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
		listTasksOutput, err := client.ListTasks(context.TODO(),
			&ecs.ListTasksInput{
				Cluster: &cluster,
			},
		)
		if err != nil {
			return nil, err
		}
		if len(listTasksOutput.TaskArns) == 0 {
			return []ecsTypes.Task{}, nil
		}
		return DescribeTask(region, cluster, listTasksOutput.TaskArns)
	}
	return DescribeTask(region, cluster, names)
}

// GetTaskInService returns tasks in the Service in the Cluster
func GetTaskInService(region string, cluster string, service string, names []string) ([]ecsTypes.Task, error) {
	client := GetClient(region)
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
	return DescribeTask(region, cluster, taskArns)
}

// filterTaskArns selects tasks specified in --name options
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

// DescribeTask returns Task list. This requires specifying task name
func DescribeTask(region string, cluster string, names []string) ([]ecsTypes.Task, error) {
	client := GetClient(region)
	describeTasksOutput, err := client.DescribeTasks(context.TODO(),
		&ecs.DescribeTasksInput{
			Cluster: &cluster,
			Tasks:   names,
		},
	)
	if err != nil {
		return []ecsTypes.Task{}, err
	}
	return describeTasksOutput.Tasks, err
}
