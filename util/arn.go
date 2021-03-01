package util

import "strings"

func ARNtoName(arn string) string {
	splited := strings.Split(arn, "/")
	return splited[len(splited)-1]
}
