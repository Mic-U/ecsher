package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Mic-U/ecsher/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

const EcsherVersion = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ecsher",
	Short:   "ecsher is the CLI tool describing AWS ECS resources.",
	Long:    `ecsher is the CLI tool describing AWS ECS resources.`,
	Version: EcsherVersion,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ecsher.toml)") // Cobra also supports local flags, which will only run
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cfgFile := filepath.Join(home, config.DefaultConfigFileName)
		if _, err := os.Stat(cfgFile); err != nil {
			_, err := os.Create(cfgFile)
			cobra.CheckErr(err)
		}
		cobra.CheckErr(err)
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// }
	err := viper.ReadInConfig()
	cobra.CheckErr(err)
	if err := viper.Unmarshal(&config.EcsherConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
