package chat

import (
	"context"
)

type Service struct {
	client Client
}

func NewService(model string) (*Service, error) {
	client, err := NewOpenAIClient(model)
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
