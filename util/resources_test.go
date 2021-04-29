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

func TestLikeInstance(t *testing.T) {
	if !LikeInstance("instance") {
		t.Fatal("instance is LikeInstance")
	} else if !LikeInstance("instances") {
		t.Fatal("instances is LikeInstance")
	}
}

func TestLikeCapacityProvider(t *testing.T) {
	if !LikeCapacityProvider("capacityprovider") {
		t.Fatal("capacityprovider is LikeInstance")
	} else if !LikeCapacityProvider("capacityproviders") {
		t.Fatal("capacityproviders is LikeInstance")
	} else if !LikeCapacityProvider("cp") {
		t.Fatal("cp is LikeInstance")
	}
}

func TestDivideTaskDefinitionArn(t *testing.T) {
	cases := []struct {
		Arn       string
		expected1 string
		expected2 string
	}{
		{
			Arn:       "arn:aws:ecs:region:account:task-definition/family:1",
			expected1: "family",
			expected2: "1",
		},
		{
			Arn:       "arn:aws:ecs:region:account:task-definition/family1",
			expected1: "family1",
			expected2: "",
		},
	}
	for _, c := range cases {
		actual1, actual2 := DivideTaskDefinitionArn(c.Arn)
		if actual1 != c.expected1 || actual2 != c.expected2 {
			t.Errorf("expect [%v, %v], got [%v, %v]", c.expected1, c.expected2, actual1, actual2)
		}
	}
}
