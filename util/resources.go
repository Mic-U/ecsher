package util

import "strings"

func LikeCluster(arg string) bool {
	return strings.Contains(strings.ToLower(arg), "cluster")
}
