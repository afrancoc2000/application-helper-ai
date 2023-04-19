package openai

import (
	"context"
	"fmt"
	"strings"
	"time"

	openAI "github.com/PullRequestInc/go-gpt3"
	"github.com/afrancoc2000/application-helper-ai/internal/config"
	gptEncoder "github.com/samber/go-gpt-3-encoder"
	azureOpenAI "github.com/sozercan/kubectl-ai/pkg/gpt3"
)

const (
	numberOfChoices = 1
	reservedTokens  = 200
	baseContext     = "You are a coding assistant for developers, you help developers create applications, you specify one by one the files needed to build an application telling the file name, the file path and the file content. You specify the file path as a valid relative path starting with a point '.'. You must return the answer as a json array, the user is a computer that needs to be able to parse your answer. You don't give explanations you don't show the commands needed to run."
)

type AIClient interface {
	QueryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error)
}

type Message struct {
	role    Role
	content string
}

func NewAIClient(appConfig config.AppConfig) (AIClient, error) {
	isChat := isChat(appConfig.OpenAIDeploymentName)
	isOpenAI := isOpenAI(appConfig.AzureOpenAIEndpoint)
	if isOpenAI {
		client := openAI.NewClient(appConfig.OpenAIAPIKey)
		if isChat {
			messages := initializeMessages()
			return &openAIChatClient{client: client, appConfig: appConfig, messages: messages}, nil
		} else {
			return &openAICompletionClient{client: client, appConfig: appConfig}, nil
		}
	} else {
		client, err := azureOpenAI.NewClient(
			appConfig.AzureOpenAIEndpoint,
			appConfig.OpenAIAPIKey,
			appConfig.OpenAIDeploymentName,
			azureOpenAI.WithTimeout(60*time.Second))
		if err != nil {
			return nil, err
		}

		if isChat {
			messages := initializeMessages()
			return &azureAIChatClient{client: client, appConfig: appConfig, messages: messages}, nil
		} else {
			return &azureAICompletionClient{client: client, appConfig: appConfig}, nil
		}
	}
}

func isChat(deployment string) bool {
	return deployment == "gpt-3.5-turbo-0301" ||
		deployment == "gpt-3.5-turbo" ||
		deployment == "gpt-4-0314" ||
		deployment == "gpt-4-32k-0314"
}

func isOpenAI(endpoint string) bool {
	return endpoint == ""
}

type openAICompletionClient struct {
	client    openAI.Client
	appConfig config.AppConfig
}

type openAIChatClient struct {
	client    openAI.Client
	appConfig config.AppConfig
	messages  []Message
}

type azureAICompletionClient struct {
	client    azureOpenAI.Client
	appConfig config.AppConfig
}

type azureAIChatClient struct {
	client    azureOpenAI.Client
	appConfig config.AppConfig
	messages  []Message
}

func calculateCompletionParams(prompts []string, appConfig config.AppConfig) (*int, *int, *strings.Builder, error) {
	choices := int(numberOfChoices)
	maxTokens, err := calculateMaxTokens(strings.Join(prompts, "\n"), appConfig.OpenAIDeploymentName, appConfig.MaxTokens)
	if err != nil {
		return nil, nil, nil, err
	}

	var prompt strings.Builder
	for _, p := range prompts {
		fmt.Fprintf(&prompt, "%s\n", p)
	}

	return &choices, maxTokens, &prompt, err
}

func calculateChatParams(messages []Message, appConfig config.AppConfig) (*int, *int, *strings.Builder, error) {
	choices := int(numberOfChoices)
	maxTokens, err := calculateMaxTokens(strings.Join(prompts, "\n"), appConfig.OpenAIDeploymentName, appConfig.MaxTokens)
	if err != nil {
		return nil, nil, nil, err
	}

	var prompt strings.Builder
	for _, p := range prompts {
		fmt.Fprintf(&prompt, "%s\n", p)
	}

	return &choices, maxTokens, &prompt, err
}

func (c *openAICompletionClient) QueryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error) {
	choices, maxTokens, prompt, err := calculateCompletionParams(prompts, c.appConfig)
	if err != nil {
		return "", err
	}

	resp, err := c.client.CompletionWithEngine(ctx, c.appConfig.OpenAIDeploymentName, openAI.CompletionRequest{
		Prompt:      []string{prompt.String()},
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           choices,
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

func (c *openAIChatClient) QueryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error) {
	choices, maxTokens, prompt, err := calculateCompletionParams(prompts, c.appConfig)
	if err != nil {
		return "", err
	}

	resp, err := c.client.ChatCompletion(ctx, openAI.ChatCompletionRequest{
		Model: c.appConfig.OpenAIDeploymentName,
		Messages: []openAI.ChatCompletionRequestMessage{
			{
				Role:    userRole,
				Content: prompt.String(),
			},
		},
		MaxTokens:   *maxTokens,
		N:           *choices,
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

func (c *azureAICompletionClient) QueryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error) {
	choices, maxTokens, prompt, err := calculateCompletionParams(prompts, c.appConfig)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Completion(ctx, azureOpenAI.CompletionRequest{
		Prompt:      []string{prompt.String()},
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           choices,
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

func (c *azureAIChatClient) QueryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error) {
	choices, maxTokens, prompt, err := calculateCompletionParams(prompts, c.appConfig)
	if err != nil {
		return "", err
	}

	resp, err := c.client.ChatCompletion(ctx, azureOpenAI.ChatCompletionRequest{
		Model: c.appConfig.OpenAIDeploymentName,
		Messages: []azureOpenAI.ChatCompletionRequestMessage{
			{
				Role:    "system",
				Content: fmt.Sprintf("%s\n%s\n", baseContext, c.appConfig.ChatContext),
			},
			{
				Role:    userRole,
				Content: prompt.String(),
			},
		},
		MaxTokens:   *maxTokens,
		N:           *choices,
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

var maxTokensMap = map[string]int{
	"code-davinci-002":   8001,
	"text-davinci-003":   4097,
	"gpt-3.5-turbo-0301": 4096,
	"gpt-3.5-turbo":      4096,
	"gpt-35-turbo-0301":  4096, // for azure
	"gpt-4-0314":         8192,
	"gpt-4-32k-0314":     8192,
}

func calculateMaxTokens(prompt string, deploymentName string, userMaxTokens int) (*int, error) {
	maxTokens := userMaxTokens
	if maxTokens == 0 {
		var ok bool
		maxTokens, ok = maxTokensMap[deploymentName]
		if !ok {
			return nil, fmt.Errorf("deploymentName %q not found in max tokens map", deploymentName)
		}
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

func convertToOpenAIMessages(messages []Message) []azureOpenAI.ChatCompletionRequestMessage  {
	openAIMessages := []azureOpenAI.ChatCompletionRequestMessage{}
	
	for _, message := range messages {
		openAIMessage := azureOpenAI.ChatCompletionRequestMessage{
			Role: message.role.String(),
			Content: message.content,
		}
		openAIMessages = append(openAIMessages, openAIMessage)
	}
	
	return openAIMessages
}

func initializeMessages() []Message {
	messages := []Message{}
	contextMessage := Message{
		role: System,
		content: baseContext,
	}
	messages = append(messages, contextMessage)
	return messages	
}