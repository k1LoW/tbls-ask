package openai

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/k1LoW/repin"
	"github.com/k1LoW/tbls/schema"
	"github.com/sashabaranov/go-openai"
)

const (
	DefaultModelChat = "gpt-4-1106-preview"
	quoteStart       = "```sql"
	quoteEnd         = "```"
)

type OpenAI struct {
	client          *openai.Client
	model           string
	promptTmpl      string
	queryPromptTmpl string
}

func New(key, model string) *OpenAI {
	return &OpenAI{
		client:          openai.NewClient(key),
		model:           model,
		promptTmpl:      defaultPromtTmpl,
		queryPromptTmpl: defaultQueryPromptTmpl,
	}
}

func (o *OpenAI) Ask(ctx context.Context, q string, s *schema.Schema) (string, error) {
	tpl, err := template.New("").Parse(o.promptTmpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, map[string]any{
		"DatabaseVersion": databaseVersion(s),
		"QuoteStart":      quoteStart,
		"QuoteEnd":        quoteEnd,
		"DDL":             generateDDLRoughly(s),
		"Question":        q,
	}); err != nil {
		return "", err
	}
	res, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       o.model,
		Temperature: 0.5, // https://community.openai.com/t/cheat-sheet-mastering-temperature-and-top-p-in-chatgpt-api-a-few-tips-and-tricks-on-controlling-the-creativity-deterministic-output-of-prompt-responses/172683
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: buf.String(),
			},
		},
	})
	if err != nil {
		return "", err
	}
	answer := res.Choices[0].Message.Content
	return answer, nil
}

func (o *OpenAI) AskQuery(ctx context.Context, q string, s *schema.Schema) (string, error) {
	tpl, err := template.New("").Parse(o.queryPromptTmpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, map[string]any{
		"DatabaseVersion": databaseVersion(s),
		"QuoteStart":      quoteStart,
		"QuoteEnd":        quoteEnd,
		"DDL":             generateDDLRoughly(s),
		"Question":        q,
	}); err != nil {
		return "", err
	}
	res, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       o.model,
		Temperature: 0.2, // https://community.openai.com/t/cheat-sheet-mastering-temperature-and-top-p-in-chatgpt-api-a-few-tips-and-tricks-on-controlling-the-creativity-deterministic-output-of-prompt-responses/172683
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: buf.String(),
			},
		},
	})
	if err != nil {
		return "", err
	}
	answer := res.Choices[0].Message.Content
	if !strings.Contains(answer, quoteStart) {
		return "", fmt.Errorf("failed to pick query from answer: %w\nanswer: %s\n", err, answer)
	}
	query := new(bytes.Buffer)
	if _, err := repin.Pick(strings.NewReader(answer), quoteStart, quoteEnd, true, query); err != nil {
		return "", fmt.Errorf("failed to pick query from answer: %w\nanswer: %s\n", err, answer)
	}
	return query.String(), nil
}

func (o *OpenAI) SetPromtTmpl(t string) {
	o.promptTmpl = t
}

func (o *OpenAI) SetQueryPromtTmpl(t string) {
	o.queryPromptTmpl = t
}

func databaseVersion(s *schema.Schema) string {
	var n string
	switch s.Driver.Name {
	case "mysql":
		n = "MySQL"
	case "sqlite":
		n = "SQLite"
	case "postgres":
		n = "PostgreSQL"
	default:
		n = s.Driver.Name
	}
	if s.Driver.DatabaseVersion != "" {
		n += " " + s.Driver.DatabaseVersion
	}
	if n == "" {
		n = "unknown"
	}
	return fmt.Sprintf("Database is %s.", n)
}
