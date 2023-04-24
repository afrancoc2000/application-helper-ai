package main

import (
	"fmt"
	"os"

	"github.com/afrancoc2000/application-helper-ai/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
