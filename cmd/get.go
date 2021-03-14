package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Mic-U/ecsher/config"
	"github.com/Mic-U/ecsher/ecs"
	util "github.com/Mic-U/ecsher/util"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get RESOURCE(cluster, service, task, definition)",
	Short: "Display resources",
	Long:  `Prints a table of important information about the specifird resources. You can filter the list using --name flag.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must specify resource")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		if util.LikeCluster(resource) {
			getCluster()
		} else if util.LikeService(resource) {
			getService()
		} else if util.LikeTask(resource) {
			getTask()
		} else if util.LikeDefinition(resource) {
			getDefinition()
		} else if util.LikeInstance(resource) {
			getInstance()
		} else {
			fmt.Printf("%s is not ECS resource\n", resource)
			os.Exit(1)
		}
	},
	Example: `  # List clusters
  ecsher get cluster
  # List services filtering by name
  esher get service -c CLUSTER_NAME --name SERVICE_NAME
  # List Tasks
  esher get task -c CLUSTER_NAME
  # List TaskDefinition families
  ecsher get definition
  ecsher get definition --prefix FAMILY_PREFIX
  # List TaskDefinition revisions in the specified family
  ecsher get definition --family FAMILY_NAME
  # List Container Instances in the specified cluster
  ecsher get instance -c CLUSTER_NAME
  `,
	Version: EcsherVersion,
}

// GetOptions used in get command
type GetOptions struct {
	Names   []string
	Cluster string
	Service string
	Region  string
	Status  string
	Prefix  string
	Family  string
}

var getOptions GetOptions

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getCmd.Flags().StringArrayVarP(&getOptions.Names, "name", "n", []string{}, "Resource name")
	getCmd.Flags().StringVarP(&getOptions.Cluster, "cluster", "c", "", "Cluster name")
	getCmd.Flags().StringVarP(&getOptions.Service, "service", "s", "", "Service name")
	getCmd.Flags().StringVarP(&getOptions.Region, "region", "r", "", "Region")
	getCmd.Flags().StringVar(&getOptions.Status, "status", "ACTIVE", "TaskDefinition status(ACTIVE or INACTIVE)")
	getCmd.Flags().StringVar(&getOptions.Prefix, "prefix", "", "TaskDefinition name prefix")
	getCmd.Flags().StringVar(&getOptions.Family, "family", "", "TaskDefinition family name")
}

func getCluster() {
	clusters, err := ecs.GetCluster(getOptions.Region, getOptions.Names)
	cobra.CheckErr(err)
	if len(clusters) == 0 {
		fmt.Println("No clusters found")
		os.Exit(1)
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "NAME \tSTATUS\tACTIVE_SERVICES\tRUNNING_TASKS\tPENDING_TASKS\tCONTAINER_INSTANCES")
	for _, cluster := range clusters {
		fmt.Fprintf(w, "%s \t%s\t%d\t%d\t%d\t%d\n",
			*cluster.ClusterName,
			*cluster.Status,
			cluster.ActiveServicesCount,
			cluster.RunningTasksCount,
			cluster.PendingTasksCount,
			cluster.RegisteredContainerInstancesCount,
		)
	}
	w.Flush()
}

func getService() {
	cluster := config.EcsherConfigManager.GetCluster(getOptions.Cluster)
	fmt.Printf("Cluster: %s\n", cluster)
	services, err := ecs.GetService(getOptions.Region, cluster, getOptions.Names)
	cobra.CheckErr(err)
	if len(services) == 0 {
		fmt.Println("No services found")
		os.Exit(1)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "NAME \tSTATUS\tLAUNCH_TYPE\tSCHEDULING_STRATEGY\tDESIRED\tRUNNING\tPENDING")
	for _, service := range services {
		fmt.Fprintf(w, "%s \t%s\t%s\t%s\t%d\t%d\t%d\n",
			*service.ServiceName,
			*service.Status,
			service.LaunchType,
			service.SchedulingStrategy,
			service.DesiredCount,
			service.RunningCount,
			service.PendingCount,
		)
	}
	w.Flush()
}

func getTask() {
	cluster := config.EcsherConfigManager.GetCluster(getOptions.Cluster)
	fmt.Printf("Cluster: %s\n", cluster)
	if getOptions.Service != "" {
		fmt.Printf("Service: %s\n", getOptions.Service)
	}
	tasks, err := ecs.GetTask(getOptions.Region, cluster, getOptions.Service, getOptions.Names)
	cobra.CheckErr(err)
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		os.Exit(1)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "NAME \tLAUNCH_TYPE \tGROUP \tCONNECTIVITY \tDESIRED_STATUS \tLAST_STATUS \tHEALTH_STATUS")
	for _, task := range tasks {
		fmt.Fprintf(w, "%s \t%s \t%s \t%s \t%s \t%s \t%s \n",
			util.ArnToName(*task.TaskArn),
			task.LaunchType,
			*task.Group,
			task.Connectivity,
			*task.DesiredStatus,
			*task.LastStatus,
			task.HealthStatus,
		)
	}
	w.Flush()
}

func getDefinition() {
	if getOptions.Family != "" {
		showTaskDefinitionRevisions()
	} else {
		showTaskDefinitionFamilies()
	}
}

func showTaskDefinitionFamilies() {
	families, err := ecs.ListFamily(getOptions.Region, getOptions.Prefix, getOptions.Status)
	cobra.CheckErr(err)
	if len(families) == 0 {
		fmt.Println("No task definitions found")
		os.Exit(1)
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "FAMILY")
	for _, family := range families {
		fmt.Fprintf(w, "%s\n", family)
	}
	w.Flush()
}

func showTaskDefinitionRevisions() {
	definitions, err := ecs.GetRevisions(getOptions.Region, getOptions.Family, getOptions.Status)
	cobra.CheckErr(err)
	if len(definitions) == 0 {
		fmt.Println("No task definitions found")
		os.Exit(1)
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "FAMILY \tREVISION")
	for _, definition := range definitions {
		family, revision := util.DivideTaskDefinitionArn(definition)
		fmt.Fprintf(w, "%s \t%s\n", family, revision)
	}
	w.Flush()
}

func getInstance() {
	cluster := config.EcsherConfigManager.GetCluster(getOptions.Cluster)
	instances, err := ecs.GetInstance(getOptions.Region, cluster, getOptions.Names)
	cobra.CheckErr(err)
	if len(instances) == 0 {
		fmt.Println("No container instances found")
		os.Exit(1)
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tEC2_INSTANCE_ID\tSTATUS\tDOCKER_VERSION\tAGENT_VERSION\tCONNECTED\tREMAINING_CPU\tREMAINING_MEMORY\tRUNNING\tPENDING")
	for _, instance := range instances {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%t\t%s\t%s\t%d\t%d\n",
			util.ArnToName(*instance.ContainerInstanceArn),
			*instance.Ec2InstanceId,
			*instance.Status,
			*instance.VersionInfo.DockerVersion,
			*instance.VersionInfo.AgentVersion,
			instance.AgentConnected,
			util.GetRemainingCpuString(instance.RemainingResources),
			util.GetRemainingMemoryString(instance.RemainingResources),
			instance.RunningTasksCount,
			instance.PendingTasksCount,
		)
	}
	w.Flush()
}
