package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// ECSClient is client instance to call AWS ECS APIs
var ECSClient *ecs.Client = nil

// GetClient is singleton function returns ECSClient
func GetClient(region string) *ecs.Client {

	if ECSClient != nil {
		return ECSClient
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	if region != "" {
		cfg.Region = *aws.String(region)
	}

	ECSClient = ecs.NewFromConfig(cfg)
	return ECSClient
}
