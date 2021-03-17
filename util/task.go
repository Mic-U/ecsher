package util

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
