package util

import (
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func AscendingSortLogs(logs []types.ServiceEvent) []types.ServiceEvent {
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
