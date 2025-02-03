package chat

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	client Client
}

func NewService(model string) (*Service, error) {
	var client Client
	var err error

	if strings.HasPrefix(model, "gpt") || strings.HasPrefix(model, "o") {
		client, err = NewOpenAIClient(model)
	} else if strings.HasPrefix(model, "gemini") {
		client, err = NewGeminiClient(model)
	} else {
		return nil, fmt.Errorf("unsupported model: %s", model)
	}

	if err != nil {
		return nil, err
	}

	return &Service{client: client}, nil
}

func (s *Service) Ask(ctx context.Context, messages []Message, queryMode bool) (string, error) {
	resp, err := s.client.Ask(ctx, messages)
	if err != nil {
		return "", err
	}

	if queryMode {
		return ExtractQuery(resp)
	}

	return resp, nil
}
