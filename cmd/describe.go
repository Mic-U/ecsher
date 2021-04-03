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
	Use:   "describe RESOURCE(cluster, service, task, definition, instance)",
	Short: "Describe detail infomation about the resource",
	Long:  `Prints detail information about the specifird resources.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must specify resource")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]

		switch {
		case util.LikeCluster(resource):
			describeCluster()
		case util.LikeService(resource):
			describeService()
		case util.LikeTask(resource):
			describeTask()
		case util.LikeDefinition(resource):
			describeDefinition()
		case util.LikeInstance(resource):
			describeInstance()
		default:
			fmt.Printf("%s is not ECS resource\n", resource)
			os.Exit(1)
		}
	},
	Example: `  # Describe Cluster
  ecsher describe cluster --name CLUSTER_NAME
  # Describe Servic
  esher describe service -c CLUSTER_NAME --name SERVICE_NAME
  # Describe Task
  ecsher describe task -c CLUSTER_NAME --name TASK_NAME
  # Describe TaskDefinition
  ecsher describe definition --family FAMILY_NAME --revision REVISION_NUMBER
  # Describe ContainerInstance
  ecsher describe instance --name CONTAINER_INSTANCE_NAME -c CLUSTER_NAME
  `,
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
	client := ecs.GetClient(describeOptions.Region, RootOptions.profile)
	clusters, err := ecs.DescribeCluster(client, []string{describeOptions.Name})
	if len(clusters) == 0 {
		fmt.Println("No cluster found")
		os.Exit(1)
	}
	cobra.CheckErr(err)
	cluster, err := yaml.Marshal(&clusters[0])
	cobra.CheckErr(err)
	fmt.Println(string(cluster))
}

func describeService() {
	if describeOptions.Name == "" {
		fmt.Println("Must specify --name")
		os.Exit(1)
	}
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster, RootOptions.profile)
	client := ecs.GetClient(describeOptions.Region, RootOptions.profile)
	services, err := ecs.DescribeService(client, cluster, []string{describeOptions.Name})
	if len(services) == 0 {
		fmt.Println("No service found")
		os.Exit(1)
	}
	cobra.CheckErr(err)
	service, err := yaml.Marshal(&services[0])
	cobra.CheckErr(err)
	fmt.Println(string(service))
}

func describeTask() {
	if describeOptions.Name == "" {
		fmt.Println("Must specify --name")
		os.Exit(1)
	}
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster, RootOptions.profile)
	client := ecs.GetClient(describeOptions.Region, RootOptions.profile)
	tasks, err := ecs.DescribeTask(client, cluster, []string{describeOptions.Name})
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		os.Exit(1)
	}
	cobra.CheckErr(err)
	task, err := yaml.Marshal(&tasks[0])
	cobra.CheckErr(err)
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
	client := ecs.GetClient(describeOptions.Region, RootOptions.profile)
	definition, err := ecs.DescribeDefinition(client, definitionName)
	cobra.CheckErr(err)
	yamlDefinition, err := yaml.Marshal(definition)
	cobra.CheckErr(err)
	fmt.Println(string(yamlDefinition))
}

func describeInstance() {
	if describeOptions.Name == "" {
		fmt.Println("Must specify --name")
		os.Exit(1)
	}
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster, RootOptions.profile)
	client := ecs.GetClient(describeOptions.Region, RootOptions.profile)
	instances, err := ecs.DescribeInstance(client, cluster, []string{describeOptions.Name})
	if len(instances) == 0 {
		fmt.Println("No container instances found")
		os.Exit(1)
	}
	cobra.CheckErr(err)

	instance, err := yaml.Marshal(&instances[0])
	cobra.CheckErr(err)
	fmt.Println(string(instance))
}
