package util

import (
	"sort"
	"time"

	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

const (
	LogGroupKey        = "awslogs-group"
	LogRegionKey       = "awslogs-region"
	LogStreamPrefixKey = "awslogs-stream-prefix"
)

type LogInformation struct {
	LogGroup  string
	LogStream string
	Region    string
}

func AscendingSortServiceLogs(logs []types.ServiceEvent) []types.ServiceEvent {
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].CreatedAt.Before(*logs[j].CreatedAt)
	})
	return logs
}

func FilterServiceLogsByTymestamp(serviceLogs []types.ServiceEvent, latestTimestamp *time.Time) []types.ServiceEvent {
	result := []types.ServiceEvent{}
	for _, l := range serviceLogs {
		if l.CreatedAt.After(*latestTimestamp) {
			result = append(result, l)
		}
	}
	return result
}

func GetLogInformation(container types.ContainerDefinition, taskID string) *LogInformation {
	containerName := *container.Name
	logGroup := container.LogConfiguration.Options[LogGroupKey]
	logStreamPrefix := container.LogConfiguration.Options[LogStreamPrefixKey]
	logStream := logStreamPrefix + "/" + containerName + "/" + taskID
	region := container.LogConfiguration.Options[LogRegionKey]
	return &LogInformation{
		LogGroup:  logGroup,
		LogStream: logStream,
		Region:    region,
	}
}

func AscendingSortTaskLogs(logs []cloudwatchTypes.OutputLogEvent) []cloudwatchTypes.OutputLogEvent {
	sort.Slice(logs, func(i, j int) bool {
		return *logs[i].Timestamp < *logs[j].Timestamp
	})
	return logs
}
