package openai

import (
	"fmt"
	"testing"

	"github.com/afrancoc2000/application-helper-ai/internal/config"
	"github.com/afrancoc2000/application-helper-ai/internal/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOpenAI(t *testing.T) {
	Convey("AIClient", t, func() {

		appConfig := config.AppConfig{
			OpenaiApiKey:         "123456",
			OpenaiDeploymentName: "gpt-4-0314",
			OpenaiDeployment:     models.Gpt4_0314,
			AzureOpenaiEndpoint:  "https://devsquad-openai-lab.openai.azure.com/",
			MaxTokens:            1000,
			SkipConfirmation:     true,
			Temperature:          0.3,
			ChatContext:          "You create html applications",
			Choices:              1,
		}

		Convey("NewAIClient Azure Completion", func() {
			appConfig.OpenaiDeploymentName = "text-davinci-003"
			appConfig.OpenaiDeployment = models.TextDavinci003

			client, err := NewAIClient(appConfig)
			So(err, ShouldBeNil)

			completionClient, ok := client.(*azureAICompletionClient)

			So(ok, ShouldEqual, true)
			So(completionClient.appConfig, ShouldResemble, appConfig)
			So(len(completionClient.prompts), ShouldEqual, 4)
		})

		Convey("NewAIClient Azure Chat", func() {
			client, err := NewAIClient(appConfig)
			So(err, ShouldBeNil)

			chatClient, ok := client.(*azureAIChatClient)

			So(ok, ShouldEqual, true)
			So(chatClient.appConfig, ShouldResemble, appConfig)
			So(len(chatClient.messages), ShouldEqual, 3)
		})

		Convey("NewAIClient OpenAI Completion", func() {
			appConfig.OpenaiDeploymentName = "text-davinci-003"
			appConfig.OpenaiDeployment = models.TextDavinci003
			appConfig.AzureOpenaiEndpoint = ""

			client, err := NewAIClient(appConfig)
			So(err, ShouldBeNil)

			completionClient, ok := client.(*openAICompletionClient)

			So(ok, ShouldEqual, true)
			So(completionClient.appConfig, ShouldResemble, appConfig)
			So(len(completionClient.prompts), ShouldEqual, 4)
		})

		Convey("NewAIClient OpenAI Chat", func() {
			appConfig.AzureOpenaiEndpoint = ""

			client, err := NewAIClient(appConfig)
			So(err, ShouldBeNil)

			chatClient, ok := client.(*openAIChatClient)

			So(ok, ShouldEqual, true)
			So(chatClient.appConfig, ShouldResemble, appConfig)
			So(len(chatClient.messages), ShouldEqual, 3)
		})

		Convey("calculateMaxTokens no user tokens", func() {
			tokens, err := calculateMaxTokens("hello", models.Gpt4_0314, 0)

			So(err, ShouldBeNil)
			So(*tokens, ShouldEqual, 7991)
		})

		Convey("calculateMaxTokens good user tokens", func() {
			tokens, err := calculateMaxTokens("hello", models.Gpt4_0314, 1000)

			So(err, ShouldBeNil)
			So(*tokens, ShouldEqual, 799)
		})

		Convey("calculateMaxTokens bad user tokens", func() {
			tokens, err := calculateMaxTokens("hello", models.Gpt4_0314, 10000)

			So(err, ShouldBeNil)
			So(*tokens, ShouldEqual, 7991)
		})

		Convey("initializeMessages", func() {
			chatContext := "You create html applications"
			messages := initializeMessages(chatContext)

			So(len(messages), ShouldEqual, 3)
			So(messages[0].Role, ShouldEqual, models.System)
			So(messages[0].Content, ShouldEqual, fmt.Sprintf("%s\n%s", baseContext, chatContext))
			So(messages[1].Role, ShouldEqual, models.User)
			So(messages[1].Content, ShouldEqual, examplePrompt)
			So(messages[2].Role, ShouldEqual, models.Assistant)
			So(messages[2].Content, ShouldEqual, `[{"fileName":"main.tf","filePath":"./","fileContent":"\n# Configure the Azure provider\nprovider \"azurerm\" {\n\tfeatures {}\n}\n\n# Create a resource group\nresource \"azurerm_resource_group\" \"aks\" {\n\tname     = var.resource_group_name\n\tlocation = var.resource_group_location\n}\n"}]`)
		})

		Convey("initializePrompts", func() {
			chatContext := "You create html applications"
			prompts := initializePrompts(chatContext)

			So(len(prompts), ShouldEqual, 4)
			So(prompts[0], ShouldEqual, baseContext)
			So(prompts[1], ShouldEqual, chatContext)
			So(prompts[2], ShouldEqual, `An example answer for the question "Create a terraform project for a resource group" would be:`)
			So(prompts[3], ShouldEqual, `[{"fileName":"main.tf","filePath":"./","fileContent":"\n# Configure the Azure provider\nprovider \"azurerm\" {\n\tfeatures {}\n}\n\n# Create a resource group\nresource \"azurerm_resource_group\" \"aks\" {\n\tname     = var.resource_group_name\n\tlocation = var.resource_group_location\n}\n"}]`)
		})

	})

}
