package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	//"github.com/rebuy-de/aws-nuke/cmd"
	"github.com/beanlearninggo/hello/cmd"
)
 
func HandleRequest(ctx context.Context) (string, error) {
	fmt.Printf("Hello, world.")
	if err := cmd.NewRootCommand().Execute(); err != nil {
		return "", err
	}
	return "Hello, world.", nil
}

func main() {
	lambda.Start(HandleRequest)
}
