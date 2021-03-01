package util

import "strings"

func LikeCluster(arg string) bool {
	return strings.Contains(strings.ToLower(arg), "cluster")
}

func LikeService(arg string) bool {
	if arg == "svc" {
		return true
	}
	return strings.Contains(strings.ToLower(arg), "service")
}
