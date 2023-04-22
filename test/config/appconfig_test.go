package models

import (
	"testing"

	"github.com/afrancoc2000/application-helper-ai/internal/config"
	"github.com/afrancoc2000/application-helper-ai/internal/models"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestAppConfig(t *testing.T) {
	Convey("AppConfig", t, func() {

		viperConfig := viper.New()

		viperConfig.Set(config.OpenaiApiKeyLabel, "123456")
		viperConfig.Set(config.OpenaiDeploymentNameLabel, "gpt-4-0314")
		viperConfig.Set(config.AzureOpenaiEndpointLabel, "https://devsquad-openai-lab.openai.azure.com/")
		viperConfig.Set(config.MaxTokensLabel, 1000)
		viperConfig.Set(config.SkipConfirmationLabel, "true")
		viperConfig.Set(config.TemperatureLabel, "0.3")
		viperConfig.Set(config.ChatContextLabel, "You create html applications")

		Convey("Initialize", func() {
			appConfig := config.AppConfig{}
			err := appConfig.Initialize(*viperConfig)

			So(err, ShouldBeNil)
			So(appConfig.OpenaiApiKey, ShouldEqual, "123456")
			So(appConfig.OpenaiDeploymentName, ShouldEqual, "gpt-4-0314")
			So(appConfig.OpenaiDeployment, ShouldEqual, models.Gpt4_0314)
			So(appConfig.AzureOpenaiEndpoint, ShouldEqual, "https://devsquad-openai-lab.openai.azure.com/")
			So(appConfig.MaxTokens, ShouldEqual, 1000)
			So(appConfig.SkipConfirmation, ShouldEqual, true)
			So(appConfig.Temperature, ShouldEqual, 0.3)
			So(appConfig.ChatContext, ShouldEqual, "You create html applications")
			So(appConfig.Choices, ShouldEqual, 1)
		})
	})

}
