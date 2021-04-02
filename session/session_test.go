package session

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockSSMPluginRunner struct{}

func (m mockSSMPluginRunner) InteractiveRun(name string, args []string) error {
	return nil
}

func TestNewSSMPluginCommand(t *testing.T) {
	NewSSMPluginCommand("ap-northeast-1")
}
func TestStart(t *testing.T) {
	cmd := SSMPluginCommand{
		runner: mockSSMPluginRunner{},
		region: "ap-northeast-1",
	}
	session := &ecsTypes.Session{
		SessionId:  aws.String("abcdefg"),
		StreamUrl:  aws.String("wss://hoge.fuga"),
		TokenValue: aws.String("aaaaaaaaaaaaaaaaaaaaaaa"),
	}
	err := cmd.Start(session)
	if err != nil {
		t.Errorf(err.Error())
	}
}
