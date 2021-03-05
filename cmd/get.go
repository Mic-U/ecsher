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
		} else {
			fmt.Printf("%s is not ECS resource\n", resource)
		}
	},
	Example: `  # List clusters
  ecsher get cluster
  # List services filtering by name
  esher get service -c CLUSTER_NAME --name SERVICE_NAME`,
}

type GetOptions struct {
	Names   []string
	Cluster string
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
	getCmd.Flags().StringVarP(&getOptions.Cluster, "cluster", "c", "default", "Cluster name")
}

func getCluster() {
	clusters, err := ecs.GetCluster(getOptions.Names)
	if err != nil {
		panic(err)
	}
	if len(clusters) == 0 {
		fmt.Println("No clusters found")
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "NAME\tSTATUS\tACTIVE_SERVICES\tRUNNING_TASKS\tPENDING_TASKS\tCONTAINER_INSTANCES")
	for _, cluster := range clusters {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%d\n",
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
	services, err := ecs.GetService(cluster, getOptions.Names)
	if err != nil {
		panic(err)
	}
	if len(services) == 0 {
		fmt.Println("No services found")
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "NAME\tSTATUS\tLAUNCH_TYPE\tSCHEDULING_STRATEGY\tDESIRED\tRUNNING\tPENDING")
	for _, service := range services {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%d\t%d\n",
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
