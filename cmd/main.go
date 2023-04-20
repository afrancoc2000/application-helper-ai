package main

import (
	"fmt"
	"os"

	"github.com/afrancoc2000/application-helper-ai/internal/cli"
	"github.com/afrancoc2000/application-helper-ai/internal/config"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	command, err := cli.NewCommand(*appConfig)
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
