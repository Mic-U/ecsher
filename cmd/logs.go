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

	"github.com/Mic-U/ecsher/cloudwatch"
	"github.com/Mic-U/ecsher/config"
	"github.com/Mic-U/ecsher/ecs"
	"github.com/Mic-U/ecsher/util"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
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
		case util.LikeTask(resource):
			getTaskLogs()
		default:
			fmt.Printf("logs command not support %s\n", resource)
			os.Exit(1)
		}
	},
}

type LogsOptions struct {
	Name      string
	Cluster   string
	Region    string
	Container string
	Watch     bool
}

var logsOptions LogsOptions

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringVarP(&logsOptions.Cluster, "cluster", "c", "", "Cluster name")
	logsCmd.Flags().StringVarP(&logsOptions.Name, "name", "n", "", "Resource name")
	logsCmd.Flags().StringVar(&logsOptions.Container, "container", "", "Container name(must specify for task log)")
	logsCmd.Flags().StringVarP(&logsOptions.Region, "region", "r", "", "Region name")
	logsCmd.Flags().BoolVarP(&logsOptions.Watch, "watch", "w", false, "Watching logs")
}

func getServiceLogs() {
	if logsOptions.Name == "" {
		fmt.Println("Must specify --name")
		os.Exit(1)
	}
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster, RootOptions.profile)
	client := ecs.GetClient(describeOptions.Region, RootOptions.profile)
	service, err := ecs.DescribeService(client, cluster, []string{logsOptions.Name})
	cobra.CheckErr(err)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tID\tMESSAGE")
	eventLogs := util.AscendingSortServiceLogs(service[0].Events)
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
		additionalLogs = util.AscendingSortServiceLogs(additionalLogs)
		for _, eventLog := range additionalLogs {
			fmt.Fprintf(w, "%s\t%s\t%s\n", eventLog.CreatedAt, *eventLog.Id, *eventLog.Message)
		}
		w.Flush()
		if len(additionalLogs) > 0 {
			latestTimestamp = additionalLogs[len(additionalLogs)-1].CreatedAt
		}
	}
}

func getTaskLogs() {
	if logsOptions.Name == "" {
		fmt.Println("Must specify --name")
		os.Exit(1)
	}
	cluster := config.EcsherConfigManager.GetCluster(describeOptions.Cluster, RootOptions.profile)
	ecsClient := ecs.GetClient(describeOptions.Region, RootOptions.profile)

	// GetTask
	tasks, err := ecs.DescribeTask(ecsClient, cluster, []string{logsOptions.Name})
	cobra.CheckErr(err)

	// GetTaskDefinition
	taskDefinitionArn := tasks[0].TaskDefinitionArn
	taskDefinition, err := ecs.DescribeDefinition(ecsClient, *taskDefinitionArn)
	cobra.CheckErr(err)

	// GetContainerDefinition
	var container ecsTypes.ContainerDefinition
	if len(taskDefinition.ContainerDefinitions) > 1 {
		if logsOptions.Container == "" {
			fmt.Println("Must specify --container when you select the task which has multiple containers")
			fmt.Println("This task has the following containers")
			for _, c := range util.ListContainerNames(taskDefinition.ContainerDefinitions) {
				fmt.Printf("- %s\n", c)
			}
			os.Exit(1)
		}
		container, err = util.ChooseContainer(taskDefinition.ContainerDefinitions, logsOptions.Container)
		cobra.CheckErr(err)
	} else {
		container = taskDefinition.ContainerDefinitions[0]
	}

	if !util.IsAwslogsLogDriver(container) {
		fmt.Println("logDriver must be awslogs")
		os.Exit(1)
	}

	logInformation := util.GetLogInformation(container, util.ArnToName(*tasks[0].TaskArn))
	cloudwatchClient := cloudwatch.GetClient(logInformation.Region)
	taskLogs, err := cloudwatch.GetTaskLog(cloudwatchClient, logInformation.LogGroup, logInformation.LogStream, nil)
	cobra.CheckErr(err)
	taskLogs = util.AscendingSortTaskLogs(taskLogs)
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tMESSAGE")
	for _, taskLog := range taskLogs {
		fmt.Fprintf(w, "%s\t%s\n", time.Unix(*taskLog.Timestamp/1000, 0), *taskLog.Message)
	}
	w.Flush()

	if logsOptions.Watch {
		latestTimeStamp := taskLogs[len(taskLogs)-1].Timestamp
		watchTaskLogs(cloudwatchClient, logInformation, latestTimeStamp, w)
	}
}

func watchTaskLogs(cloudwatchClient cloudwatchlogs.GetLogEventsAPIClient, logInformation *util.LogInformation, latestTimeStamp *int64, w *tabwriter.Writer) {
	for {
		time.Sleep(time.Second * 5)
		taskLogs, err := cloudwatch.GetTaskLog(cloudwatchClient, logInformation.LogGroup, logInformation.LogStream, latestTimeStamp)
		cobra.CheckErr(err)
		taskLogs = util.AscendingSortTaskLogs(taskLogs)
		for _, taskLog := range taskLogs {
			fmt.Fprintf(w, "%s\t%s\n", time.Unix(*taskLog.Timestamp/1000, 0), *taskLog.Message)
		}
		w.Flush()
		if len(taskLogs) > 0 {
			latestTimeStamp = taskLogs[len(taskLogs)-1].Timestamp
		}
	}
}
