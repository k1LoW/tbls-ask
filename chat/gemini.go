package chat

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiClient(model string) (*GeminiClient, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &GeminiClient{
		client: client,
		model:  client.GenerativeModel(model),
	}, nil
}

func (c *GeminiClient) Ask(ctx context.Context, messages []Message) (string, error) {
	chat := c.model.StartChat()

	// Convert messages to Gemini format
	history := make([]*genai.Content, len(messages))
	for i, msg := range messages {
		role := msg.Role
		switch role {
		case "system":
			role = "user"
		case "assistant":
			role = "model"
		}

		history[i] = &genai.Content{
			Parts: []genai.Part{
				genai.Text(msg.Content),
			},
			Role: role,
		}
	}

	chat.History = history

	// Send the last message
	lastMsg := messages[len(messages)-1]
	resp, err := chat.SendMessage(ctx, genai.Text(lastMsg.Content))
	if err != nil {
		return "", fmt.Errorf("gemini api error: %w", err)
	}

	// Extract response
	var answer string
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			for _, part := range candidate.Content.Parts {
				if part != nil {
					answer = fmt.Sprintf("%s", part)
				}
			}
		}
	}
	return answer, nil
}
