package main

import (
	//"github.com/rebuy-de/aws-nuke/cmd"
	"github.com/beanlearninggo/hello/cmd"
)

type NukeEvent struct {
	AccountId   string `json:"account_id"`
	IAMUsername string `json:"iam_username"`
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
}

func HandleRequest(event NukeEvent) (string, error) {
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
	HandleRequest(
		NukeEvent{
			AccountId:   "356077346614",
			IAMUsername: "nukesurvivoradmin",
			AccessKey:   "AKIAVFZ65D43LMRX6MSD",
			SecretKey:   "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",
		},
	)
}
