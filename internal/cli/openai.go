package cli

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	openai "github.com/PullRequestInc/go-gpt3"
	gptEncoder "github.com/samber/go-gpt-3-encoder"
	azureopenai "github.com/sozercan/kubectl-ai/pkg/gpt3"
)

const (
	userRole        = "user"
	numberOfChoices = 1
	baseContext     = "You are a coding assistant for developers, you help developers create applications, you specify one by one the files needed to build an application telling the file name, the file path and the file content. you don't give explanations you don't show the commands needed to run"
)

type AIClient interface {
	queryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error)
}

func newAIClient() (AIClient, error) {
	if azureOpenAIEndpoint == nil || *azureOpenAIEndpoint == "" {
		client := openai.NewClient(*openAIAPIKey)
		if isChat() {
			return &openAIChatClient{client: client}, nil
		} else {
			return &openAICompletionClient{client: client}, nil
		}
	} else {
		re := regexp.MustCompile(`^[a-zA-Z0-9]+([_-]?[a-zA-Z0-9]+)*$`)
		if !re.MatchString(*openAIDeploymentName) {
			err := errors.New("azure openai deployment can only include alphanumeric characters, '_,-', and can't end with '_' or '-'")
			return nil, err
		}

		client, err := azureopenai.NewClient(*azureOpenAIEndpoint, *openAIAPIKey, *openAIDeploymentName)
		if err != nil {
			return nil, err
		}

		if isChat() {
			return &azureAIChatClient{client: client}, nil
		} else {
			return &azureAICompletionClient{client: client}, nil
		}
	}
}

func isChat() bool {
	return *openAIDeploymentName == "gpt-3.5-turbo-0301" ||
		*openAIDeploymentName == "gpt-3.5-turbo" ||
		*openAIDeploymentName == "gpt-4-0314" ||
		*openAIDeploymentName == "gpt-4-32k-0314"
}

type openAICompletionClient struct {
	client openai.Client
}

type openAIChatClient struct {
	client openai.Client
}

type azureAICompletionClient struct {
	client azureopenai.Client
}

type azureAIChatClient struct {
	client azureopenai.Client
}

func calculateParams(prompts []string, deploymentName string) (*float32, *int, *int, *strings.Builder, error) {
	temp := float32(*temperature)
	choices := int(numberOfChoices)
	maxTokens, err := calculateMaxTokens(prompts, deploymentName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var prompt strings.Builder
	fmt.Fprintf(&prompt, "%s\n", baseContext)
	fmt.Fprintf(&prompt, "%s\n", *chatContext)
	for _, p := range prompts {
		fmt.Fprintf(&prompt, "%s\n", p)
	}

	return &temp, &choices, maxTokens, &prompt, err
}

func (c *openAICompletionClient) queryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error) {
	temp, choices, maxTokens, prompt, err := calculateParams(prompts, deploymentName)
	if err != nil {
		return "", err
	}

	resp, err := c.client.CompletionWithEngine(ctx, *openAIDeploymentName, openai.CompletionRequest{
		Prompt:      []string{prompt.String()},
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           choices,
		Temperature: temp,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Text, nil
}

func (c *openAIChatClient) queryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error) {
	temp, choices, maxTokens, prompt, err := calculateParams(prompts, deploymentName)
	if err != nil {
		return "", err
	}

	resp, err := c.client.ChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: *openAIDeploymentName,
		Messages: []openai.ChatCompletionRequestMessage{
			{
				Role:    userRole,
				Content: prompt.String(),
			},
		},
		MaxTokens:   *maxTokens,
		N:           *choices,
		Temperature: temp,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *azureAICompletionClient) queryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error) {
	temp, choices, maxTokens, prompt, err := calculateParams(prompts, deploymentName)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Completion(ctx, azureopenai.CompletionRequest{
		Prompt:      []string{prompt.String()},
		MaxTokens:   maxTokens,
		Echo:        false,
		N:           choices,
		Temperature: temp,
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("expected choices to be 1 but received: %d", len(resp.Choices))
	}

	return resp.Choices[0].Text, nil
}

func (c *azureAIChatClient) queryOpenAI(ctx context.Context, prompts []string, deploymentName string) (string, error) {
	temp, choices, maxTokens, prompt, err := calculateParams(prompts, deploymentName)
	if err != nil {
		return "", err
	}

	resp, err := c.client.ChatCompletion(ctx, azureopenai.ChatCompletionRequest{
		Model: *openAIDeploymentName,
		Messages: []azureopenai.ChatCompletionRequestMessage{
			{
				Role:    userRole,
				Content: prompt.String(),
			},
		},
		MaxTokens:   *maxTokens,
		N:           *choices,
		Temperature: temp,
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

func calculateMaxTokens(prompts []string, deploymentName string) (*int, error) {
	var maxTokensFinal int
	if *maxTokens == 0 {
		var ok bool
		maxTokensFinal, ok = maxTokensMap[deploymentName]
		if !ok {
			return nil, fmt.Errorf("deploymentName %q not found in max tokens map", deploymentName)
		}
	} else {
		maxTokensFinal = *maxTokens
	}

	encoder, err := gptEncoder.NewEncoder()
	if err != nil {
		return nil, err
	}

	// start at 100 since the encoder at times doesn't get it exactly correct
	totalTokens := 100
	for _, prompt := range prompts {
		tokens, err := encoder.Encode(prompt)
		if err != nil {
			return nil, err
		}
		totalTokens += len(tokens)
	}

	remainingTokens := maxTokensFinal - totalTokens
	return &remainingTokens, nil
}
