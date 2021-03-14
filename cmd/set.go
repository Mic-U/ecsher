// Package cmd defins each commands(get, set and so on)
package cmd

import (
	"errors"
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
	Use:   "set RESOURCE(cluster)",
	Short: "Set cluster persistently",
	Long: `It's very annoying to enter the cluster name every time.
When you use set command, ecsher remembers it`,
	Example: ` # Set cluster name
  ecsher set cluster --name CLUSTER_NAME`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must specify resource")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		if util.LikeCluster(resource) {
			setCluster()
		} else {
			fmt.Printf("ecsher set does not support %s currently\n", resource)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().StringVarP(&setOptions.Name, "name", "n", "", "Resource name")
}

func setCluster() {
	err := config.EcsherConfigManager.SetCluster(setOptions.Name)
	cobra.CheckErr(err)
	fmt.Printf("Cluster: %s\n", setOptions.Name)

}
