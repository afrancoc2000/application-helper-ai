package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/afrancoc2000/application-helper-ai/internal/config"
	fileSystem "github.com/afrancoc2000/application-helper-ai/internal/file_system"
	"github.com/afrancoc2000/application-helper-ai/internal/models"
	"github.com/afrancoc2000/application-helper-ai/internal/openai"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	version = "0.0.1"

	apply      = "Apply"
	doNotApply = "Don't apply"
	makeBetter = "Add to the query"
)

type Command struct {
	client      openai.AIClient
	fileFactory fileSystem.FileFactory
	appConfig   config.AppConfig
}

func NewCommand(appConfig config.AppConfig) (*Command, error) {
	flag.Parse()

	if appConfig.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("Please provide an OpenAI key.")
	}

	client, err := openai.NewAIClient(appConfig)
	if err != nil {
		return nil, err
	}

	fileFactory := fileSystem.NewFileFactory()

	return &Command{
		client:      client,
		fileFactory: &fileFactory,
		appConfig:   appConfig,
	}, nil
}

func (c *Command) CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "application-ai",
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("prompt must be provided")
			}

			err := c.run(args)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (c *Command) run(args []string) error {
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

		text := fmt.Sprintf(
			"These are the files that would be created. Do you want to apply them? or add something to the query?\n%s",
			queryResult)
		files, err = models.AppFileFromString(text)
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

func (c *Command) userActionPrompt() (string, error) {
	// if require confirmation is not set, immediately return apply
	if !c.appConfig.SkipConfirmation {
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
	for index, file := range files {
		fmt.Printf("%d. File: %s/%s:\n", index, file.Path, file.Name)
		fmt.Printf("%s\n", file.Content)
		fmt.Printf("\n")
	}
}
