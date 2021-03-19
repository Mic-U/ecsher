package cloudwatch

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type mockTaskLogClient struct {
	output cloudwatchlogs.GetLogEventsOutput
}

func (m mockTaskLogClient) GetLogEvents(ctx context.Context, params *cloudwatchlogs.GetLogEventsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogEventsOutput, error) {
	if *params.LogGroupName == "" || *params.LogStreamName == "" {
		return nil, errors.New("LogGroupName and LogStreamName is needful")
	}
	if params.StartTime == nil {
		return &m.output, nil
	}

	result := []cloudwatchTypes.OutputLogEvent{}
	for _, e := range m.output.Events {
		if *e.Timestamp >= *params.StartTime {
			result = append(result, e)
		}
	}
	return &cloudwatchlogs.GetLogEventsOutput{
		Events: result,
	}, nil
}

func TestGetTaskLog(t *testing.T) {
	cases := []struct {
		Client    mockTaskLogClient
		LogGroup  string
		LogStream string
		StartTime *int64
	}{
		{
			Client: mockTaskLogClient{
				output: cloudwatchlogs.GetLogEventsOutput{
					Events: []cloudwatchTypes.OutputLogEvent{
						{
							Timestamp: aws.Int64(100),
							Message:   aws.String("test"),
						},
					},
				},
			},
			LogGroup:  "test",
			LogStream: "test",
			StartTime: aws.Int64(10),
		},
		{
			Client: mockTaskLogClient{
				output: cloudwatchlogs.GetLogEventsOutput{
					Events: []cloudwatchTypes.OutputLogEvent{
						{
							Timestamp: aws.Int64(100),
							Message:   aws.String("test"),
						},
					},
				},
			},
			LogGroup:  "",
			LogStream: "test",
			StartTime: aws.Int64(10),
		},
		{
			Client: mockTaskLogClient{
				output: cloudwatchlogs.GetLogEventsOutput{
					Events: []cloudwatchTypes.OutputLogEvent{
						{
							Timestamp: aws.Int64(100),
							Message:   aws.String("test"),
						},
					},
				},
			},
			LogGroup:  "test",
			LogStream: "",
			StartTime: aws.Int64(10),
		},
		{
			Client: mockTaskLogClient{
				output: cloudwatchlogs.GetLogEventsOutput{
					Events: []cloudwatchTypes.OutputLogEvent{
						{
							Timestamp: aws.Int64(100),
							Message:   aws.String("test"),
						},
					},
				},
			},
			LogGroup:  "test",
			LogStream: "test",
			StartTime: nil,
		},
		{
			Client: mockTaskLogClient{
				output: cloudwatchlogs.GetLogEventsOutput{
					Events: []cloudwatchTypes.OutputLogEvent{
						{
							Timestamp: aws.Int64(100),
							Message:   aws.String("test"),
						},
					},
				},
			},
			LogGroup:  "test",
			LogStream: "test",
			StartTime: aws.Int64(100),
		},
	}

	for _, c := range cases {
		result, err := GetTaskLog(c.Client, c.LogGroup, c.LogStream, c.StartTime)
		if c.LogGroup == "" || c.LogStream == "" {
			if err == nil {
				t.Fatal("should return err")
			}
			continue
		}
		if err != nil {
			t.Fatal("err should not occur")
		}
		if c.StartTime == nil {
			continue
		}
		if *c.StartTime == 10 {
			if len(result) != 1 {
				t.Fatalf("expected is 1, but actual is %v", len(result))
			}
		}
		if *c.StartTime == 100 {
			if len(result) != 0 {
				t.Fatalf("expected is 0, but actual is %v", len(result))
			}
		}
	}
}
