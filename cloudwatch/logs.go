package cloudwatch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

var TaskLogLimit int32 = 100

func GetTaskLog(client cloudwatchlogs.GetLogEventsAPIClient, logGroup string, logStream string, latestTimeStamp *int64) (output []cloudwatchTypes.OutputLogEvent, err error) {
	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
		Limit:         aws.Int32(TaskLogLimit),
	}

	if latestTimeStamp != nil {
		startTime := *latestTimeStamp + 1 // Start
		input.StartTime = aws.Int64(startTime)
	}
	logs, err := client.GetLogEvents(context.TODO(), input)
	if err != nil {
		return []cloudwatchTypes.OutputLogEvent{}, err
	}
	output = logs.Events
	return output, nil
}
