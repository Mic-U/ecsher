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
	"os"
	"text/tabwriter"
	"time"

	"github.com/Mic-U/ecsher/config"
	"github.com/Mic-U/ecsher/ecs"
	"github.com/Mic-U/ecsher/util"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs RESOURCE(service)",
	Short: "Prints event logs",
	Long:  `Prints event logs`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("must specify resource")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		switch {
		case util.LikeService(resource):
			getServiceLogs()
		default:
			fmt.Printf("logs command not support %s\n", resource)
			os.Exit(1)
		}
	},
}

type LogsOptions struct {
	Name    string
	Cluster string
	Region  string
	Watch   bool
}

var logsOptions LogsOptions

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringVarP(&logsOptions.Cluster, "cluster", "c", "", "Cluster name")
	logsCmd.Flags().StringVarP(&logsOptions.Name, "name", "n", "", "Resource name")
	logsCmd.Flags().StringVarP(&logsOptions.Region, "region", "r", "", "Region name")
	logsCmd.Flags().BoolVarP(&logsOptions.Watch, "watch", "w", false, "Watching logs")
}

func getServiceLogs() {
	if logsOptions.Name == "" {
		fmt.Println("Must specify --name")
		os.Exit(1)
	}
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster)
	client := ecs.GetClient(describeOptions.Region)
	service, err := ecs.DescribeService(client, cluster, []string{logsOptions.Name})
	cobra.CheckErr(err)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tID\tMESSAGE")
	eventLogs := util.AscendingSortLogs(service[0].Events)
	for _, eventLog := range eventLogs {
		fmt.Fprintf(w, "%s\t%s\t%s\n", eventLog.CreatedAt, *eventLog.Id, *eventLog.Message)
	}
	w.Flush()

	if logsOptions.Watch {
		latestTimestamp := eventLogs[len(eventLogs)-1].CreatedAt
		watchServiceLogs(client, cluster, latestTimestamp, w)
	}
}

func watchServiceLogs(client ecs.ECSServiceClient, cluster string, latestTimestamp *time.Time, w *tabwriter.Writer) {
	for {
		time.Sleep(time.Second * 5)
		service, err := ecs.DescribeService(client, cluster, []string{logsOptions.Name})
		cobra.CheckErr(err)
		additionalLogs := service[0].Events
		additionalLogs = util.FilterServiceLogsByTymestamp(additionalLogs, latestTimestamp)
		additionalLogs = util.AscendingSortLogs(additionalLogs)
		for _, eventLog := range additionalLogs {
			fmt.Fprintf(w, "%s\t%s\t%s\n", eventLog.CreatedAt, *eventLog.Id, *eventLog.Message)
		}
		w.Flush()
		if len(additionalLogs) > 0 {
			latestTimestamp = additionalLogs[len(additionalLogs)-1].CreatedAt
		}
	}
}
