package config

import "github.com/kelseyhightower/envconfig"

type AppConfig struct {
	OpenAIAPIKey         string  `required:"true"  envconfig:"OPENAI_API_KEY"`
	OpenAIDeploymentName string  `required:"false" envconfig:"OPENAI_DEPLOYMENT_NAME" default:"text-davinci-003"`
	MaxTokens            int     `required:"false" envconfig:"MAX_TOKENS" default:"0"`
	AzureOpenAIEndpoint  string  `required:"false" envconfig:"AZURE_OPENAI_ENDPOINT" default:""`
	SkipConfirmation     bool    `required:"false" envconfig:"SKIP_CONFIRMATION" default:"false"`
	Temperature          float32 `required:"false" envconfig:"TEMPERATURE" default:"0"`
	ChatContext          string  `required:"false" envconfig:"OPENAI_CHAT_CONTEXT" default:""`
}

func NewAppConfig() (*AppConfig, error) {
	config := &AppConfig{}
	err := envconfig.Process("", config)
	return config, err
}
