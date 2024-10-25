package gemini

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Gemini struct {
	chatClient          *genai.Client
	model           string
}

func NewClient(ctx context.Context, model string) (*Gemini, error) {
	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}
	return &Gemini{
		chatClient:          client,
		model:           model,
	}, nil
}

func (g *Gemini) ChatCompletionRequest(ctx context.Context, p string) (string, error) {
	model := g.chatClient.GenerativeModel(g.model)
	resp, err := model.GenerateContent(ctx, genai.Text(p))
	if err != nil {
		return "", err
	}
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
