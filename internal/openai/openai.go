package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openAI "github.com/PullRequestInc/go-gpt3"
	"github.com/afrancoc2000/application-helper-ai/internal/config"
	"github.com/afrancoc2000/application-helper-ai/internal/models"
	gptEncoder "github.com/samber/go-gpt-3-encoder"
	azureOpenAI "github.com/sozercan/kubectl-ai/pkg/gpt3"
)

const (
	numberOfChoices      = 1
	reservedTokens       = 200
	baseContext          = "You are a coding assistant for developers, you help developers create applications, you specify one by one the files needed to build an application telling the file name, the file path and the file content. You specify the file path as a valid relative path starting with a point '.'. You must return the answer as a json array, the user is a computer that needs to be able to parse your answer. You don't give explanations you don't show the commands needed to run."
	examplePrompt        = "Create a terraform project for a resource group"
	exampleAnswerName    = "main.tf"
	exampleAnswerPath    = "./"
	exampleAnswerContent = `
	# Configure the Azure provider  
	provider "azurerm" {  
		features {}  
	}  
	  
	# Create a resource group  
	resource "azurerm_resource_group" "aks" {  
		name     = var.resource_group_name  
		location = var.resource_group_location  
	}  		
	`
)

type AIClient interface {
	QueryOpenAI(ctx context.Context, prompt string) (string, error)
}

func NewAIClient(appConfig config.AppConfig) (AIClient, error) {
	isChat := isChat(appConfig.OpenAIDeployment)
	isOpenAI := isOpenAI(appConfig.AzureOpenAIEndpoint)
	if isOpenAI {
		client := openAI.NewClient(appConfig.OpenAIAPIKey)
		if isChat {
			messages := initializeMessages()
			return &openAIChatClient{client: client, appConfig: appConfig, messages: messages}, nil
		} else {
			prompts := initializePrompts()
			return &openAICompletionClient{client: client, appConfig: appConfig, prompts: prompts}, nil
		}
	} else {
		client, err := azureOpenAI.NewClient(
			appConfig.AzureOpenAIEndpoint,
			appConfig.OpenAIAPIKey,
			appConfig.OpenAIDeployment.String(),
			azureOpenAI.WithTimeout(60*time.Second))
		if err != nil {
			return nil, err
		}

		if isChat {
			messages := initializeMessages()
			return &azureAIChatClient{client: client, appConfig: appConfig, messages: messages}, nil
		} else {
			prompts := initializePrompts()
			return &azureAICompletionClient{client: client, appConfig: appConfig, prompts: prompts}, nil
		}
	}
}

func isChat(deployment models.Deployment) bool {
	return deployment.IsChat()
}

func isOpenAI(endpoint string) bool {
	return endpoint == ""
}

type openAICompletionClient struct {
	client    openAI.Client
	appConfig config.AppConfig
	prompts   []string
}

type openAIChatClient struct {
	client    openAI.Client
	appConfig config.AppConfig
	messages  []models.Message
}

type azureAICompletionClient struct {
	client    azureOpenAI.Client
	appConfig config.AppConfig
	prompts   []string
}

type azureAIChatClient struct {
	client    azureOpenAI.Client
	appConfig config.AppConfig
	messages  []models.Message
}

func calculateCompletionParams(prompts []string, appConfig config.AppConfig) (*int, error) {
	return calculateMaxTokens(strings.Join(prompts, "\n"), appConfig.OpenAIDeployment, appConfig.MaxTokens)
}

func calculateChatParams(messages []models.Message, appConfig config.AppConfig) (*int, error) {
	prompts, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}
	return calculateMaxTokens(string(prompts), appConfig.OpenAIDeployment, appConfig.MaxTokens)
}

func (c *openAICompletionClient) QueryOpenAI(ctx context.Context, prompt string) (string, error) {
	c.prompts = append(c.prompts, prompt)
	maxTokens, err := calculateCompletionParams(c.prompts, c.appConfig)
	if err != nil {
		return "", err
	}

	resp, err := c.client.CompletionWithEngine(ctx, c.appConfig.OpenAIDeployment.String(), openAI.CompletionRequest{
		Prompt:      c.prompts,
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           &c.appConfig.Choices,
		Temperature: &c.appConfig.Temperature,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Text, nil
}

func (c *openAIChatClient) QueryOpenAI(ctx context.Context, prompt string) (string, error) {
	message := models.Message{
		Role:    models.User,
		Content: prompt,
	}
	c.messages = append(c.messages, message)
	maxTokens, err := calculateChatParams(c.messages, c.appConfig)
	if err != nil {
		return "", err
	}

	resp, err := c.client.ChatCompletion(ctx, openAI.ChatCompletionRequest{
		Model:       c.appConfig.OpenAIDeployment.String(),
		Messages:    models.ConvertToOpenAIMessages(c.messages),
		MaxTokens:   *maxTokens,
		N:           *&c.appConfig.Choices,
		Temperature: &c.appConfig.Temperature,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("expected choices to be 1 but received: %d", len(resp.Choices))
	}

	message = models.Message{
		Role:    models.Assistant,
		Content: resp.Choices[0].Message.Content,
	}
	c.messages = append(c.messages, message)

	return resp.Choices[0].Message.Content, nil
}

func (c *azureAICompletionClient) QueryOpenAI(ctx context.Context, prompt string) (string, error) {
	c.prompts = append(c.prompts, prompt)
	maxTokens, err := calculateCompletionParams(c.prompts, c.appConfig)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Completion(ctx, azureOpenAI.CompletionRequest{
		Prompt:      c.prompts,
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           &c.appConfig.Choices,
		Temperature: &c.appConfig.Temperature,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Text, nil
}

func (c *azureAIChatClient) QueryOpenAI(ctx context.Context, prompt string) (string, error) {
	message := models.Message{
		Role:    models.User,
		Content: prompt,
	}
	c.messages = append(c.messages, message)
	maxTokens, err := calculateChatParams(c.messages, c.appConfig)
	if err != nil {
		return "", err
	}

	resp, err := c.client.ChatCompletion(ctx, azureOpenAI.ChatCompletionRequest{
		Model:       c.appConfig.OpenAIDeployment.String(),
		Messages:    models.ConvertToAzureOpenAIMessages(c.messages),
		MaxTokens:   *maxTokens,
		N:           *&c.appConfig.Choices,
		Temperature: &c.appConfig.Temperature,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Message.Content, nil
}

func calculateMaxTokens(prompt string, deployment models.Deployment, userMaxTokens int) (*int, error) {
	maxTokens := userMaxTokens
	if maxTokens == 0 {
		maxTokens = deployment.MaxTokens()
	}

	encoder, err := gptEncoder.NewEncoder()
	if err != nil {
		return nil, err
	}

	totalTokens := reservedTokens
	tokens, err := encoder.Encode(prompt)
	if err != nil {
		return nil, err
	}
	totalTokens += len(tokens)

	remainingTokens := maxTokens - totalTokens
	return &remainingTokens, nil
}

func initializeMessages() []models.Message {
	messages := []models.Message{}
	contextMessage := models.Message{
		Role:    models.System,
		Content: baseContext,
	}
	messages = append(messages, contextMessage)

	examplePromptMessage := models.Message{
		Role:    models.User,
		Content: examplePrompt,
	}
	messages = append(messages, examplePromptMessage)

	exampleAnswer := models.AppFile{
		Name:    exampleAnswerName,
		Path:    exampleAnswerPath,
		Content: exampleAnswerContent,
	}
	jsonContent, _ := json.Marshal(exampleAnswer)

	exampleAnswerMessage := models.Message{
		Role:    models.Assistant,
		Content: string(jsonContent),
	}
	messages = append(messages, exampleAnswerMessage)

	return messages
}

func initializePrompts() []string {
	prompts := []string{baseContext}
	return prompts
}
