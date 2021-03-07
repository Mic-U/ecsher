package cmd

import (
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
	Use:   "get",
	Short: "Display resources",
	Long: `Prints a table of important information about the specifird resources. You can filter the list using --name flag.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Specify resource. Please see 'esher help get'")
			return
		}
		resource := args[0]
		if util.LikeCluster(resource) {
			getCluster()
		} else if util.LikeService(resource) {
			getService()
		} else if util.LikeTask(resource) {
			getTask()
		} else {
			fmt.Printf("%s is not ECS resource\n", resource)
		}
	},
	Example: `  # List clusters
  ecsher get cluster
  # List services filtering by name
  esher get service -c CLUSTER_NAME --name SERVICE_NAME`,
}

// GetOptions used in get command
type GetOptions struct {
	Names   []string
	Cluster string
	Service string
	Region  string
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
}

func getCluster() {
	clusters, err := ecs.GetCluster(getOptions.Region, getOptions.Names)
	if err != nil {
		panic(err)
	}
	if len(clusters) == 0 {
		fmt.Println("No clusters found")
		return
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
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
	if err != nil {
		panic(err)
	}
	if len(services) == 0 {
		fmt.Println("No services found")
		return
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
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
	if err != nil {
		panic(err)
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
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
