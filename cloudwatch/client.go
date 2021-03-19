package cloudwatch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

var LogsClient *cloudwatchlogs.Client = nil

func GetClient(region string) *cloudwatchlogs.Client {
	if LogsClient != nil {
		return LogsClient
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	if region != "" {
		cfg.Region = *aws.String(region)
	}
	LogsClient = cloudwatchlogs.NewFromConfig(cfg)
	return LogsClient
}
