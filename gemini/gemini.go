package gemini

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/k1LoW/tbls-ask/templates"
	"github.com/k1LoW/tbls/schema"
	"google.golang.org/api/option"
)

const (
	DefaultModelChat = "gemini-pro"
	quoteStart       = "```sql"
	quoteEnd         = "```"
)

type Gemini struct {
	client          *genai.Client
	model           string
	promptTmpl      string
	queryPromptTmpl string
}

func New(key, model string) *Gemini {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(key))
	if err != nil {
		return nil
	}
	return &Gemini{
		client:          client,
		model:           model,
		promptTmpl:      templates.DefaultPromtTmpl,
		queryPromptTmpl: templates.DefaultQueryPromptTmpl,
	}
}

func (g *Gemini) Ask(ctx context.Context, q string, s *schema.Schema) (string, error) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not set")
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", err
	}
	defer client.Close()

	tpl, err := template.New("").Parse(templates.DefaultPromtTmpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, map[string]any{
		"DatabaseVersion": templates.DatabaseVersion(s),
		"QuoteStart":      "```sql",
		"QuoteEnd":        "```",
		"DDL":             templates.GenerateDDLRoughly(s),
		"Question":        q,
	}); err != nil {
		return "", err
	}

	model := client.GenerativeModel(g.model)
	resp, err := model.GenerateContent(ctx, genai.Text(buf.String()))
	if err != nil {
		return "", err
	}
	answer := extractResponse(resp)
	return answer, nil
}

func (g *Gemini) AskQuery(ctx context.Context, q string, s *schema.Schema) (string, error) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not set")
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return "", err
	}
	defer client.Close()

	tpl, err := template.New("").Parse(templates.DefaultQueryPromptTmpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, map[string]any{
		"DatabaseVersion": templates.DatabaseVersion(s),
		"QuoteStart":      "```sql",
		"QuoteEnd":        "```",
		"DDL":             templates.GenerateDDLRoughly(s),
		"Question":        q,
	}); err != nil {
		return "", err
	}

	model := client.GenerativeModel(g.model)
	resp, err := model.GenerateContent(ctx, genai.Text(buf.String()))
	if err != nil {
		return "", err
	}
	answer := extractResponse(resp)
	return answer, nil
}

func extractResponse(resp *genai.GenerateContentResponse) string {
	response := ""
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			for _, part := range candidate.Content.Parts {
				if part != nil {
					response = fmt.Sprintf("%s", part)
				}
			}
		}
	}
	return response
}
