package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetCluster(t *testing.T) {
	result := EcsherConfigManager.GetCluster("test")
	if result != "test" {
		t.Fatal("Cluster name should be 'test'")
	}
	viper.Set("cluster", "test2")
	result = EcsherConfigManager.GetCluster("test")
	if result != "test2" {
		t.Fatal("Cluster name should be 'test'")
	}
}
