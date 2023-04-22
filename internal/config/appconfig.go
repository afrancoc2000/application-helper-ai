package config

import (
	"github.com/afrancoc2000/application-helper-ai/internal/models"
	"github.com/spf13/viper"
)

const (
	OpenaiApiKeyLabel         = "openaiApiKey"
	OpenaiDeploymentNameLabel = "openaiDeploymentName"
	OpenaiDeploymentLabel     = "openaiDeployment"
	MaxTokensLabel            = "maxTokens"
	AzureOpenaiEndpointLabel  = "azureOpenaiEndpoint"
	SkipConfirmationLabel     = "skipConfirmation"
	TemperatureLabel          = "temperature"
	ChatContextLabel          = "chatContext"
	choices                   = 1
)

type AppConfig struct {
	OpenaiApiKey         string
	OpenaiDeploymentName string
	OpenaiDeployment     models.Deployment
	MaxTokens            int
	AzureOpenaiEndpoint  string
	SkipConfirmation     bool
	Temperature          float32
	ChatContext          string
	Choices              int
}

func (c *AppConfig) Initialize(viperConfig viper.Viper) error {

	c.OpenaiApiKey = viperConfig.GetString(OpenaiApiKeyLabel)
	c.OpenaiDeploymentName = viperConfig.GetString(OpenaiDeploymentNameLabel)
	c.MaxTokens = viperConfig.GetInt(MaxTokensLabel)
	c.AzureOpenaiEndpoint = viperConfig.GetString(AzureOpenaiEndpointLabel)
	c.SkipConfirmation = viperConfig.GetBool(SkipConfirmationLabel)
	c.Temperature = float32(viperConfig.GetFloat64(TemperatureLabel))
	c.ChatContext = viperConfig.GetString(ChatContextLabel)

	deployment, err := models.DeploymentFromName(c.OpenaiDeploymentName)
	if err != nil {
		return err
	}
	c.OpenaiDeployment = deployment
	c.Choices = choices

	return nil
}
