package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/walles/env"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	version = "0.0.7"

	apply      = "Apply"
	dontApply  = "Don't apply"
	makeBetter = "Add to the query"
)

var (
	kubernetesConfigFlags *genericclioptions.ConfigFlags

	openAIDeploymentName = flag.String("openai-deployment-name", env.GetOr("OPENAI_DEPLOYMENT_NAME", env.String, "text-davinci-003"), "The deployment name used for the model in OpenAI service.")
	maxTokens            = flag.Int("max-tokens", env.GetOr("MAX_TOKENS", strconv.Atoi, 0), "The max token will overwrite the max tokens in the max tokens map.")
	openAIAPIKey         = flag.String("openai-api-key", env.GetOr("OPENAI_API_KEY", env.String, ""), "The API key for the OpenAI service. This is required.")
	azureOpenAIEndpoint  = flag.String("azure-openai-endpoint", env.GetOr("AZURE_OPENAI_ENDPOINT", env.String, ""), "The endpoint for Azure OpenAI service. If provided, Azure OpenAI service will be used instead of OpenAI service.")
	requireConfirmation  = flag.Bool("require-confirmation", env.GetOr("REQUIRE_CONFIRMATION", strconv.ParseBool, true), "Whether to require confirmation before executing the command. Defaults to true.")
	temperature          = flag.Float64("temperature", env.GetOr("TEMPERATURE", env.WithBitSize(strconv.ParseFloat, 64), 0.0), "The temperature to use for the model. Range is between 0 and 1. Set closer to 0 if your want output to be more deterministic but less creative. Defaults to 0.0.")
	chatContext          = flag.String("openai-chat-context", env.GetOr("OPENAI_CHAT_CONTEXT", env.String, ""), "The text context for the OpenAI service to know what kind of app to generate.")
)

func InitAndExecute() {
	flag.Parse()

	if *openAIAPIKey == "" {
		fmt.Println("Please provide an OpenAI key.")
		os.Exit(1)
	}

	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "application-ai",
		Version:      version,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("prompt must be provided")
			}

			err := run(args)
			if err != nil {
				return err
			}

			return nil
		},
	}

	kubernetesConfigFlags = genericclioptions.NewConfigFlags(false)
	kubernetesConfigFlags.AddFlags(cmd.Flags())
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	return cmd
}

func run(args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	client, err := newAIClient()
	if err != nil {
		return err
	}

	var action, queryResult string
	for action != apply {
		args = append(args, action)
		queryResult, err = client.queryOpenAI(ctx, args, *openAIDeploymentName)
		if err != nil {
			return err
		}

		text := fmt.Sprintf(
			"These are the files that would be created. Do you want to apply them? or add something to the query?\n%s",
			queryResult)
		fmt.Println(text)

		action, err := userActionPrompt()
		if err != nil {
			return err
		}

		if action == dontApply {
			return nil
		}
	}
	return applyManifest(queryResult)
}

func userActionPrompt() (string, error) {
	// if require confirmation is not set, immediately return apply
	if !*requireConfirmation {
		return apply, nil
	}

	var result string
	var err error
	items := []string{apply, dontApply}
	label := fmt.Sprintf("Would you like to apply this? [%s/%s/%s]", makeBetter, apply, dontApply)

	prompt := promptui.SelectWithAdd{
		Label:    label,
		Items:    items,
		AddLabel: makeBetter,
	}
	_, result, err = prompt.Run()
	if err != nil {
		return dontApply, err
	}

	return result, nil
}