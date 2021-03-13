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

func TestLikeTask(t *testing.T) {
	if !LikeTask("task") {
		t.Fatal("task is LikeTask")
	} else if !LikeTask("tasks") {
		t.Fatal("tasks is LikeTask")
	} else if LikeTask("taskdef") {
		t.Fatal("taskdef is not LikeTask")
	}
}

func TestLikeDefinition(t *testing.T) {
	if !LikeDefinition("taskdef") {
		t.Fatal("taskdef is LikeDefinition")
	} else if !LikeDefinition("definition") {
		t.Fatal("definition is LikeDefinition")
	} else if !LikeDefinition("definitions") {
		t.Fatal("definitions is LikeDefinition")
	}
}
