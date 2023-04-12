package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	//"github.com/rebuy-de/aws-nuke/cmd"
	"github.com/beanlearninggo/hello/cmd"
)

type NukeEvent struct {
	AccountId   string `json:"account_id"`
	IAMUsername string `json:"iam_username"`
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
}

func HandleRequest(ctx context.Context, event NukeEvent) (string, error) {
	var err error

	command := cmd.NewRootCommand()
	command.SetArgs([]string{
		"--account-id", event.AccountId,
		"--iam-username", event.IAMUsername,
		"--access-key-id", event.AccessKey,
		"--secret-access-key", event.SecretKey,
	})
	if err = command.Execute(); err != nil {
		return "", err
	}
	return "Nuke complete", nil
}

func main() {
	lambda.Start(HandleRequest)
}
