package ecs

import "testing"

func TestGetClient(t *testing.T) {
	// Should not occurs any errors
	GetClient("ap-northeast-1", "default")
	GetClient("ap-northeast-1", "default")
}
