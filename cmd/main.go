package main

import (
	"fmt"
	"os"

	"github.com/afrancoc2000/application-helper-ai/internal/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	command, err := cli.NewCommand()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = command.CreateCommand().Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
