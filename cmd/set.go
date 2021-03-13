// Package cmd defins each commands(get, set and so on)
package cmd

import (
	"fmt"
	"os"

	"github.com/Mic-U/ecsher/config"
	util "github.com/Mic-U/ecsher/util"
	"github.com/spf13/cobra"
)

// SetOptions used in set commands
type SetOptions struct {
	Name string
}

var setOptions SetOptions

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set cluster persistently",
	Long: `It's very annoying to enter the cluster name every time.
When you use set command, ecsher remembers it`,
	Example: ` # Set cluster name
  ecsher set cluster --name CLUSTER_NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Specify resource. Please see 'ecsher help set'")
			return
		}
		resource := args[0]
		if util.LikeCluster(resource) {
			setCluster()
		}
		fmt.Println("set called")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setCmd.Flags().StringVarP(&setOptions.Name, "name", "n", "", "Resource name")
}

func setCluster() {
	err := config.EcsherConfigManager.SetCluster(setOptions.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("Cluster: %s\n", setOptions.Name)
	}
}
