package openai

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	chatClient *openai.Client
	model      string
}

func NewClient(model string) (*OpenAI, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is not set")
	}
	return &OpenAI{
		chatClient: openai.NewClient(key),
		model:      model,
	}, nil
}

func (o *OpenAI) ChatCompletionRequest(ctx context.Context, p string) (string, error) {
	res, err := o.chatClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       o.model,
		Temperature: 0.5, // https://community.openai.com/t/cheat-sheet-mastering-temperature-and-top-p-in-chatgpt-api-a-few-tips-and-tricks-on-controlling-the-creativity-deterministic-output-of-prompt-responses/172683
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: p,
			},
		},
	})
	if err != nil {
		return "", err
	}
	answer := res.Choices[0].Message.Content
	return answer, nil
}
