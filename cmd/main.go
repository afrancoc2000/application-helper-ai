package main

import (
	"github.com/afrancoc2000/application-helper-ai/internal/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	cli.InitAndExecute()
}
