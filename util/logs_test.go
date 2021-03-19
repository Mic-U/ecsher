package util

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func TestAscendingSortServiceLogs(t *testing.T) {
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
	result := AscendingSortServiceLogs(testLogs)
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

func TestGetLogInformation(t *testing.T) {
	testConfig := types.ContainerDefinition{
		Name: aws.String("test"),
		LogConfiguration: &types.LogConfiguration{
			LogDriver: "awslogs",
			Options:   map[string]string{"awslogs-group": "group", "awslogs-region": "region", "awslogs-stream-prefix": "teststream"},
		},
	}

	result := GetLogInformation(testConfig, "123456789012")
	if result.LogGroup != "group" {
		t.Errorf("expected is group, but actual is %s", result.LogGroup)
	}

	if result.Region != "region" {
		t.Errorf("expected is region, but actual is %s", result.Region)
	}

	if result.LogStream != "teststream/test/123456789012" {
		t.Errorf("expected is teststream/test/123456789012, but actual is %s", result.LogStream)
	}
}

func TestAscendingTaskLogs(t *testing.T) {
	testLogs := []cloudwatchTypes.OutputLogEvent{
		{
			Timestamp: aws.Int64(100),
		},
		{
			Timestamp: aws.Int64(101),
		},
		{
			Timestamp: aws.Int64(99),
		},
	}
	result := AscendingSortTaskLogs(testLogs)
	for i, r := range result {
		if i == 0 {
			continue
		}
		if *result[i-1].Timestamp >= (*r.Timestamp) {
			t.Fatal("Should sort by ascending order")
		}
	}
}
