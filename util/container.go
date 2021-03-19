package util

import (
	"fmt"

	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func ListContainerNames(containers []ecsTypes.ContainerDefinition) []string {
	result := []string{}
	for _, container := range containers {
		result = append(result, *container.Name)
	}
	return result
}

func ChooseContainer(containers []ecsTypes.ContainerDefinition, name string) (ecsTypes.ContainerDefinition, error) {
	for _, container := range containers {
		if *container.Name == name {
			return container, nil
		}
	}
	return ecsTypes.ContainerDefinition{}, fmt.Errorf("%s is not found in containers", name)
}
