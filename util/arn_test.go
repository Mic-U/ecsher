package util

import (
	"testing"
)

func TestClusterARNtoName(t *testing.T) {
	clusterARN := "arn:aws:ecs:us-east-1:123456789012:cluster/my-cluster"
	result := ARNtoName(clusterARN)
	if result != "my-cluster" {
		t.Fatal("ClusterARN was not translated to name correctedly")
	}
}

func TestServiceARNtoName(t *testing.T) {
	shortServiceARN := "arn:aws:ecs:us-east-1:123456789012:service/my-service"
	result := ARNtoName(shortServiceARN)
	if result != "my-service" {
		t.Fatal("ShortServiceARN was not translated to name correctedly")
	}

	longServiceARN := "arn:aws:ecs:us-east-1:123456789012:my-cluster/my-service"
	result = ARNtoName(longServiceARN)
	if result != "my-service" {
		t.Fatal("ShortServiceARN was not translated to name correctedly")
	}
}
