package client

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	

	"github.com/k1LoW/repin"
)

const (
	quoteStart       = "```sql"
	quoteEnd         = "```"
)

type LLMAgent interface {
	ChatCompletionRequest(ctx context.Context, prompt string) (string, error)
}

type Client struct {
	Agent  LLMAgent
	Querymode   bool
}

func (c *Client) Ask(ctx context.Context, p string) (string, error) {
	resp, err := c.Agent.ChatCompletionRequest(ctx, p)
	if err != nil {
		return "", err
	}
	if c.Querymode {
		resp, err = extractQuery(resp)
		if err != nil {
			return "", err
		}
	}
	return resp, nil
}

func extractQuery(resp string) (string, error) {
	if !strings.Contains(resp, quoteStart) {
		return "", fmt.Errorf("failed to pick query from answer: %s", resp)
	}
	if !strings.HasSuffix(resp, "\n") {
		resp += "\n"
	}
	src := strings.NewReader(resp)
	query := new(bytes.Buffer)
	if _, err := repin.Pick(src, quoteStart, quoteEnd, true, query); err != nil {
		return "", fmt.Errorf("failed to pick query from answer: %w\nanswer: %s", err, resp)
	}
	return query.String(), nil
}
