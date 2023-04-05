package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rebuy-de/aws-nuke/cmd"
)

func handler(ctx context.Context) error {
	fmt.Printf("Hello, world.")
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(-1)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
