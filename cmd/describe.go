package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

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
			return errors.New("requires a resource")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]

		if util.LikeCluster(resource) {
			describeCluster()
		} else if util.LikeService(resource) {
			describeService()
		} else if util.LikeTask(resource) {
			describeTask()
		} else if util.LikeDefinition(resource) {
			describeDefinition()
		}
	},
}

// DescribeOptions used in describe command
type DescribeOptions struct {
	Name     string
	Cluster  string
	Region   string
	Family   string
	Revision int
}

var describeOptions DescribeOptions

func init() {
	rootCmd.AddCommand(describeCmd)

	describeCmd.Flags().StringVarP(&describeOptions.Name, "name", "n", "", "Resource name")
	describeCmd.Flags().StringVarP(&describeOptions.Cluster, "cluster", "c", "", "Cluster name")
	describeCmd.Flags().StringVarP(&describeOptions.Region, "region", "r", "", "Region")
	describeCmd.Flags().StringVar(&describeOptions.Family, "family", "", "TaskDefinition family name")
	describeCmd.Flags().IntVar(&describeOptions.Revision, "revision", 0, "TaskDefinition revision")
}

func describeCluster() {
	if describeOptions.Name == "" {
		fmt.Println("Must specify --name")
		os.Exit(1)
	}
	clusters, err := ecs.DescribeCluster(describeOptions.Region, []string{describeOptions.Name})
	if len(clusters) == 0 {
		fmt.Println("No cluster found")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cluster, err := yaml.Marshal(&clusters[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(cluster))
}

func describeService() {
	if describeOptions.Name == "" {
		fmt.Println("Must specify --name")
		os.Exit(1)
	}
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster)
	services, err := ecs.DescribeService(describeOptions.Region, cluster, []string{describeOptions.Name})
	if len(services) == 0 {
		fmt.Println("No service found")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	service, err := yaml.Marshal(&services[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(service))
}

func describeTask() {
	if describeOptions.Name == "" {
		fmt.Println("Must specify --name")
		os.Exit(1)
	}
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster)
	tasks, err := ecs.DescribeTask(describeOptions.Region, cluster, []string{describeOptions.Name})
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	task, err := yaml.Marshal(&tasks[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(task))
}

func describeDefinition() {
	if describeOptions.Family == "" {
		fmt.Println("Must specify --family")
		os.Exit(1)
	}

	if describeOptions.Revision < 1 {
		fmt.Println("Must specify Positive number in --revision")
		os.Exit(1)
	}

	definitionName := describeOptions.Family + ":" + strconv.Itoa(describeOptions.Revision)
	definition, err := ecs.DescribeDefinition(describeOptions.Region, definitionName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	yamlDefinition, err := yaml.Marshal(definition)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(yamlDefinition))
}
