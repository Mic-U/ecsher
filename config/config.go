package config

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type EcsherConfigStruct struct {
	Cluster string
}

type ConfigNames struct {
	Cluster string
}

type EcsherConfigManagerStruct struct {
	configNames ConfigNames
}

var EcsherConfig = &EcsherConfigStruct{}

var configNames = &ConfigNames{
	Cluster: "cluster",
}

var DefaultConfigFileName string = ".ecsher.toml"

var EcsherConfigManager = &EcsherConfigManagerStruct{
	configNames: *configNames,
}

func (m EcsherConfigManagerStruct) SetCluster(cluster string, profile string) error {
	key := profile + "." + m.configNames.Cluster
	viper.Set(key, cluster)
	return m.saveConfig()
}

func (m EcsherConfigManagerStruct) GetCluster(optionCluster string, profile string) string {
	key := profile + "." + m.configNames.Cluster
	cfgFileCluster, ok := viper.Get(key).(string)
	if !ok {
		return optionCluster
	}
	if optionCluster == "" && cfgFileCluster != "" {
		return cfgFileCluster
	}
	return optionCluster
}

func (m EcsherConfigManagerStruct) saveConfig() error {
	if err := viper.WriteConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			home, err := homedir.Dir()
			if err == nil {
				err = viper.WriteConfigAs(filepath.Join(home, DefaultConfigFileName))
			}
			return err
		}
	}
	return nil
}
