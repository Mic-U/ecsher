package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetCluster(t *testing.T) {
	result := EcsherConfigManager.GetCluster("test1", "default")
	if result != "test1" {
		t.Fatal("Cluster name should be 'test1'")
	}
	viper.Set("hoge.cluster", "test2")
	result = EcsherConfigManager.GetCluster("", "hoge")
	if result != "test2" {
		t.Fatal("Cluster name should be 'test2'")
	}
	result = EcsherConfigManager.GetCluster("test3", "hoge")
	if result != "test3" {
		t.Fatal("Cluster name should be 'test3'")
	}
}

func TestSecCluster(t *testing.T) {
	result := EcsherConfigManager.SetCluster("test1", "hoge")
	if result != nil {
		t.Fatalf("expect no error, got %v", result)
	}
	cluster := viper.Get("hoge.cluster")
	if cluster != "test1" {
		t.Fatalf("expect is test1, acutual is %v", cluster)
	}

}
