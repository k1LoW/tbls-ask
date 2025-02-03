package chat

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
	model  string
}

func newAzureClient(model string) (*OpenAIClient, error) {
	key := os.Getenv("AZURE_OPENAI_KEY")
	if key == "" {
		return nil, fmt.Errorf("AZURE_OPENAI_KEY is not set")
	}
	c := openai.DefaultAzureConfig(key, os.Getenv("AZURE_OPENAI_ENDPOINT"))

	if os.Getenv("AZURE_OPENAI_API_VERSION") != "" {
		c.APIVersion = os.Getenv("AZURE_OPENAI_API_VERSION")
	}

	if os.Getenv("AZURE_OPENAI_MODEL") != "" {
		c.AzureModelMapperFunc = func(model string) string {
			return os.Getenv("AZURE_OPENAI_MODEL")
		}
	}
	return &OpenAIClient{
		client: openai.NewClientWithConfig(c),
		model:  model,
	}, nil
}

func NewOpenAIClient(model string) (*OpenAIClient, error) {
	if os.Getenv("AZURE_OPENAI_ENDPOINT") != "" {
		return newAzureClient(model)
	}

	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}
	return &OpenAIClient{
		client: openai.NewClient(key),
		model:  model,
	}, nil
}

func (c *OpenAIClient) Ask(ctx context.Context, messages []Message) (string, error) {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       c.model,
		Temperature: 0.5,
		Messages:    openaiMessages,
	})
	if err != nil {
		return "", fmt.Errorf("openai api error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
