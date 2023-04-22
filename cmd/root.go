package cmd

import (
	"fmt"
	fileSystem "github.com/afrancoc2000/application-helper-ai/internal/file_system"

	"github.com/afrancoc2000/application-helper-ai/internal/appai"
	"github.com/afrancoc2000/application-helper-ai/internal/config"
	"github.com/afrancoc2000/application-helper-ai/internal/openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "1.0.0"

var appConfig = config.AppConfig{}
var viperConfig = *viper.New()

var RootCmd = &cobra.Command{
	Use:   "application-ai",
	Short: "Application AI is an application generator that uses OpenAI",
	Long: `A command line app that receives a natural language instruction
		and simulates a conversation with OpenAI and in result, generates 
		the files defined`,
	Version:      version,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return fmt.Errorf("prompt must be provided")
		}

		err := appConfig.Initialize(viperConfig)
		if err != nil {
			return err
		}

		client, err := openai.NewAIClient(appConfig)
		if err != nil {
			return err
		}

		fileFactory := fileSystem.NewFileFactory()

		generator, err := appai.NewGenerator(appConfig, client, fileFactory)
		if err != nil {
			return err
		}

		err = generator.Run(args)
		return err
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringP(
		config.OpenaiApiKeyLabel,
		"k",
		"",
		"The API key for the OpenAI service. This is required.")

	RootCmd.PersistentFlags().StringP(
		config.OpenaiDeploymentNameLabel,
		"d",
		"text-davinci-003",
		"The deployment name in OpenAI/Azure for the model")

	RootCmd.PersistentFlags().IntP(
		config.MaxTokensLabel,
		"m",
		0,
		"The max token will overwrite the max tokens in the max tokens map.")

	RootCmd.PersistentFlags().StringP(
		config.AzureOpenaiEndpointLabel,
		"e",
		"",
		"The endpoint for Azure OpenAI service. If provided, Azure OpenAI service will be used instead of OpenAI service.")

	RootCmd.PersistentFlags().BoolP(
		config.SkipConfirmationLabel,
		"s",
		false,
		"Whether to skip confirmation before creating the files. Defaults to false.")

	RootCmd.PersistentFlags().Float32P(
		config.TemperatureLabel,
		"t",
		0,
		"The temperature to use for the model. Range is between 0 and 1. Set closer to 0 if your want output to be more deterministic but less creative. Defaults to 0.0.")

	RootCmd.PersistentFlags().StringP(
		config.ChatContextLabel,
		"c",
		"",
		"The text context for the OpenAI service to know what kind of app to generate.")
}

func initConfig() {
	viperConfig.SetEnvPrefix("")

	viperConfig.BindEnv(config.OpenaiApiKeyLabel, "OPENAI_API_KEY")
	viperConfig.BindEnv(config.OpenaiDeploymentNameLabel, "OPENAI_DEPLOYMENT_NAME")
	viperConfig.BindEnv(config.MaxTokensLabel, "MAX_TOKENS")
	viperConfig.BindEnv(config.AzureOpenaiEndpointLabel, "AZURE_OPENAI_ENDPOINT")
	viperConfig.BindEnv(config.SkipConfirmationLabel, "SKIP_CONFIRMATION")
	viperConfig.BindEnv(config.TemperatureLabel, "TEMPERATURE")
	viperConfig.BindEnv(config.ChatContextLabel, "CHAT_CONTEXT")

	viperConfig.BindPFlag(config.OpenaiApiKeyLabel, RootCmd.Flags().Lookup(config.OpenaiApiKeyLabel))
	viperConfig.BindPFlag(config.OpenaiDeploymentNameLabel, RootCmd.Flags().Lookup(config.OpenaiDeploymentNameLabel))
	viperConfig.BindPFlag(config.MaxTokensLabel, RootCmd.Flags().Lookup(config.MaxTokensLabel))
	viperConfig.BindPFlag(config.AzureOpenaiEndpointLabel, RootCmd.Flags().Lookup(config.AzureOpenaiEndpointLabel))
	viperConfig.BindPFlag(config.SkipConfirmationLabel, RootCmd.Flags().Lookup(config.SkipConfirmationLabel))
	viperConfig.BindPFlag(config.TemperatureLabel, RootCmd.Flags().Lookup(config.TemperatureLabel))
	viperConfig.BindPFlag(config.ChatContextLabel, RootCmd.Flags().Lookup(config.ChatContextLabel))

}
