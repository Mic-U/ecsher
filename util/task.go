package util

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func FilterTasksByNames(taskArns []string, names []string) []string {
	if len(names) == 0 {
		return taskArns
	}
	filteredTasks := []string{}
	for _, taskArn := range taskArns {
		taskName := ArnToName(taskArn)
		for _, name := range names {
			if taskName == ArnToName(name) {
				filteredTasks = append(filteredTasks, taskArn)
			}
		}
	}
	return filteredTasks
}

func GetCapacityProviderName(task ecsTypes.Task) string {
	if task.CapacityProviderName == nil {
		return ""
	} else {
		return *task.CapacityProviderName
	}
}
