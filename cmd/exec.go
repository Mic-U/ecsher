/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/Mic-U/ecsher/config"
	"github.com/Mic-U/ecsher/ecs"
	"github.com/Mic-U/ecsher/session"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec TASK_NAME COMMAND",
	Short: "Execution command in the container",
	Long:  `Execution command in the container`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must specify TASK_NAME and COMMAND")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		task := args[0]
		command := args[1]
		exec(task, command)
	},
}

type ExecOptions struct {
	Cluster     string
	Region      string
	Container   string
	Interactive bool
}

var execOptions ExecOptions

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&execOptions.Cluster, "cluster", "c", "", "Cluster Name")
	execCmd.Flags().StringVarP(&execOptions.Region, "region", "r", "", "Region")
	execCmd.Flags().StringVar(&execOptions.Container, "container", "", "Container Name")
	execCmd.Flags().BoolVarP(&execOptions.Interactive, "interactive", "i", false, "Interactive mode")
}

func exec(taskName string, command string) {
	region := execOptions.Region
	cluster := config.EcsherConfigManager.GetCluster(getOptions.Cluster)
	fmt.Printf("Cluster: %s\n", cluster)
	client := ecs.GetClient(region)
	output, err := ecs.ExecuteCommand(client, cluster, &ecs.ExecuteCmmandParams{
		Task:        taskName,
		Container:   execOptions.Container,
		Command:     command,
		Interactive: execOptions.Interactive,
	})
	cobra.CheckErr(err)
	cmd := session.NewSSMPluginCommand(region)
	err = cmd.Start(output.Session)
	cobra.CheckErr(err)
}
