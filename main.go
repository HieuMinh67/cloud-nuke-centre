package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	//"github.com/rebuy-de/aws-nuke/cmd"
	"github.com/beanlearninggo/hello/cmd"
)

type NukeEvent struct {
	AccountId string `json:"account_id"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

func HandleRequest(ctx context.Context, event NukeEvent) (string, error) {
	fmt.Printf("Hello, world.")
	var err error

	command := cmd.NewRootCommand()
	command.SetArgs([]string{
		"--account-id", event.AccountId,
		"--access-key-id", event.AccessKey,
		"--secret-access-key", event.SecretKey,
	})
	if err = command.Execute(); err != nil {
		return "", err
	}
	return "Hello, world.", nil
}

func main() {
	lambda.Start(HandleRequest)
}
