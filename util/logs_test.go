package util

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func TestAscendingSortLogs(t *testing.T) {
	time1, _ := time.Parse(time.RFC3339, "2021-03-01T03:00:00Z")
	time2, _ := time.Parse(time.RFC3339, "2021-03-01T03:00:01Z")
	time3, _ := time.Parse(time.RFC3339, "2021-03-01T03:00:02Z")
	testLogs := []types.ServiceEvent{
		{
			CreatedAt: &time3,
		},
		{
			CreatedAt: &time2,
		},
		{
			CreatedAt: &time1,
		},
	}
	result := AscendingSortLogs(testLogs)
	for i, r := range result {
		if i == 0 {
			continue
		}
		if result[i-1].CreatedAt.After(*r.CreatedAt) {
			t.Fatal("Should sort by ascending order")
		}
	}
}

func TestFilterServiceLogsByTymestamp(t *testing.T) {
	time1, _ := time.Parse(time.RFC3339, "2021-03-01T03:00:00Z")
	time2, _ := time.Parse(time.RFC3339, "2021-03-01T03:00:01Z")
	time3, _ := time.Parse(time.RFC3339, "2021-03-01T03:00:02Z")
	latestTimeStamp, _ := time.Parse(time.RFC3339, "2021-03-01T03:00:01Z")
	testLogs := []types.ServiceEvent{
		{
			CreatedAt: &time3,
		},
		{
			CreatedAt: &time2,
		},
		{
			CreatedAt: &time1,
		},
	}
	result := FilterServiceLogsByTymestamp(testLogs, &latestTimeStamp)
	if len(result) != 1 {
		t.Fatal("Should filtering logs")
	}
}
