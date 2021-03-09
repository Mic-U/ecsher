package cmd

import (
	"errors"
	"fmt"

	"github.com/Mic-U/ecsher/config"
	"github.com/Mic-U/ecsher/ecs"
	"github.com/Mic-U/ecsher/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe detail infomation about the resource",
	Long:  `Prints detail information about the specifird resources.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a resource")
		}
		if describeOptions.Name == "" {
			return errors.New("name is not specified. Please specify name via --name option")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Specify resource. Please see 'esher help get'")
			return
		}
		resource := args[0]

		if util.LikeCluster(resource) {
			describeCluster()
		} else if util.LikeService(resource) {
			describeService()
		} else if util.LikeTask(resource) {
			describeTask()
		}
	},
}

// DescribeOptions used in describe command
type DescribeOptions struct {
	Name    string
	Cluster string
	Region  string
}

var describeOptions DescribeOptions

func init() {
	rootCmd.AddCommand(describeCmd)

	describeCmd.Flags().StringVarP(&describeOptions.Name, "name", "n", "", "Resource name")
	describeCmd.Flags().StringVarP(&describeOptions.Cluster, "cluster", "c", "", "Cluster name")
	describeCmd.Flags().StringVarP(&describeOptions.Region, "region", "r", "", "Region")
}

func describeCluster() {
	clusters, err := ecs.DescribeCluster(describeOptions.Region, []string{describeOptions.Name})
	if err != nil {
		panic(err)
	}
	cluster, err := yaml.Marshal(&clusters[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cluster))
}

func describeService() {
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster)
	services, err := ecs.DescribeService(describeOptions.Region, cluster, []string{describeOptions.Name})
	if err != nil {
		panic(err)
	}
	service, err := yaml.Marshal(&services[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(service))
}

func describeTask() {
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster)
	tasks, err := ecs.DescribeTask(describeOptions.Region, cluster, []string{describeOptions.Name})
	if err != nil {
		panic(err)
	}
	task, err := yaml.Marshal(&tasks[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(task))
}
