package appai

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/afrancoc2000/application-helper-ai/internal/config"
	fileSystem "github.com/afrancoc2000/application-helper-ai/internal/file_system"
	"github.com/afrancoc2000/application-helper-ai/internal/models"
	"github.com/afrancoc2000/application-helper-ai/internal/openai"
	"github.com/manifoldco/promptui"
)

const (
	apply      = "Apply"
	doNotApply = "Don't apply"
	makeBetter = "Add to the query"
)

type Generator struct {
	appConfig   config.AppConfig
	client      openai.AIClient
	fileFactory fileSystem.FileFactory
}

func NewGenerator(appConfig config.AppConfig, client openai.AIClient, fileFactory fileSystem.FileFactory) (*Generator, error) {

	return &Generator{
		appConfig:   appConfig,
		client:      client,
		fileFactory: fileFactory,
	}, nil
}

func (c *Generator) Run(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	prompt := args[0]
	var action, queryResult string
	var err error
	var files []models.AppFile
	for action != apply {

		queryResult, err = c.client.QueryOpenAI(ctx, prompt)
		if err != nil {
			return err
		}

		files, err = models.AppFileFromString(queryResult)
		if err != nil {
			return err
		}

		printQueryResults(files)

		action, err = c.userActionPrompt()
		if err != nil {
			return err
		}

		if action == doNotApply {
			return nil
		}
		prompt = action
	}
	return c.fileFactory.CreateFiles(files)
}

func (c *Generator) userActionPrompt() (string, error) {
	// if skip confirmation is set, immediately return apply
	if c.appConfig.SkipConfirmation {
		return apply, nil
	}

	var result string
	var err error
	items := []string{apply, doNotApply}
	label := fmt.Sprintf("Would you like to apply this? [%s/%s/%s]", makeBetter, apply, doNotApply)

	prompt := promptui.SelectWithAdd{
		Label:    label,
		Items:    items,
		AddLabel: makeBetter,
	}
	_, result, err = prompt.Run()
	if err != nil {
		return doNotApply, err
	}

	return result, nil
}

func printQueryResults(files []models.AppFile) {
	fmt.Println("These are the files that would be created. Do you want to apply them? or add something to the query?")
	for index, file := range files {
		fmt.Printf("%d. File: %s%s:\n", index+1, file.Path, file.Name)
		fmt.Printf("%s\n", file.Content)
		fmt.Printf("\n")
	}
}
