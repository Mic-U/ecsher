package util

import "strings"

func FilterTasksByNames(taskArns []string, names []string) []string {
	filteredTasks := []string{}
	for _, taskArn := range taskArns {
		for _, name := range names {
			if strings.Contains(taskArn, name) {
				filteredTasks = append(filteredTasks, taskArn)
			}
		}
	}
	return filteredTasks
}
