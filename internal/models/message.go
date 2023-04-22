package models

import (
	openAI "github.com/PullRequestInc/go-gpt3"
	azureOpenAI "github.com/sozercan/kubectl-ai/pkg/gpt3"
)

type Message struct {
	Role    Role
	Content string
}

func ConvertToOpenAIMessages(messages []Message) []openAI.ChatCompletionRequestMessage {
	openAIMessages := []openAI.ChatCompletionRequestMessage{}

	for _, message := range messages {
		openAIMessage := openAI.ChatCompletionRequestMessage{
			Role:    message.Role.String(),
			Content: message.Content,
		}
		openAIMessages = append(openAIMessages, openAIMessage)
	}

	return openAIMessages
}

func ConvertToAzureOpenAIMessages(messages []Message) []azureOpenAI.ChatCompletionRequestMessage {
	openAIMessages := []azureOpenAI.ChatCompletionRequestMessage{}

	for _, message := range messages {
		openAIMessage := azureOpenAI.ChatCompletionRequestMessage{
			Role:    message.Role.String(),
			Content: message.Content,
		}
		openAIMessages = append(openAIMessages, openAIMessage)
	}

	return openAIMessages
}
