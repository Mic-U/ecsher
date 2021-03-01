package ecs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

var ECSClient *ecs.Client = nil

func GetClient() *ecs.Client {

	if ECSClient != nil {
		return ECSClient
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	ECSClient = ecs.NewFromConfig(cfg)
	return ECSClient
}
