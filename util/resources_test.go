package util

import (
	"testing"
)

func TestLikeCluster(t *testing.T) {
	if !LikeCluster("cluster") {
		t.Fatal("cluster is LikeCluster")
	} else if !LikeCluster("clusters") {
		t.Fatal("clusters is LikeCluster")
	}
}
