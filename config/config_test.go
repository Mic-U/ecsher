package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetCluster(t *testing.T) {
	result := EcsherConfigManager.GetCluster("test1")
	if result != "test1" {
		t.Fatal("Cluster name should be 'test1'")
	}
	viper.Set("cluster", "test2")
	result = EcsherConfigManager.GetCluster("")
	if result != "test2" {
		t.Fatal("Cluster name should be 'test2'")
	}
	result = EcsherConfigManager.GetCluster("test3")
	if result != "test3" {
		t.Fatal("Cluster name should be 'test3'")
	}
}

func TestSecCluster(t *testing.T) {
	result := EcsherConfigManager.SetCluster("test1")
	if result != nil {
		t.Fatalf("expect no error, got %v", result)
	}
}
