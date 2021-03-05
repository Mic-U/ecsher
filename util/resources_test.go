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

func TestLikeService(t *testing.T) {
	if !LikeService("service") {
		t.Fatal("service is LikeService")
	} else if !LikeService("services") {
		t.Fatal("services is LikeService")
	} else if !LikeService("svc") {
		t.Fatal("svc is LikeService")
	}
}
