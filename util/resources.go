package util

import "strings"

const (
	TaskDefinitionAlias = "taskdef"
	ServiceAlias        = "svc"
)

func LikeCluster(arg string) bool {
	return strings.Contains(strings.ToLower(arg), "cluster")
}

func LikeService(arg string) bool {
	if arg == ServiceAlias {
		return true
	}
	return strings.Contains(strings.ToLower(arg), "service")
}

func LikeTask(arg string) bool {
	if arg == TaskDefinitionAlias {
		return false
	}
	return strings.Contains(strings.ToLower(arg), "task")
}

func LikeDefinition(arg string) bool {
	if arg == TaskDefinitionAlias {
		return true
	}
	return strings.Contains(strings.ToLower(arg), "definition")
}

func LikeInstance(arg string) bool {
	return strings.Contains(strings.ToLower(arg), "instance")
}

func ArnToName(arn string) string {
	splited := strings.Split(arn, "/")
	return splited[len(splited)-1]
}

func DivideTaskDefinitionArn(taskDefinitionArn string) (string, string) {
	name := ArnToName(taskDefinitionArn)
	splited := strings.Split(name, ":")
	if len(splited) <= 1 {
		return splited[0], ""
	}
	return splited[0], splited[1]
}
