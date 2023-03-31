package main

import (
	"fmt"
	"os"

	"github.com/rebuy-de/aws-nuke/cmd"
)

func main() {
	fmt.Printf(("Hello, world."))
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(-1)
	}

}
