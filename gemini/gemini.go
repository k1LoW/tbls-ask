package gemini

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/datasource"
	"github.com/k1LoW/tbls-ask/templates"
	"google.golang.org/api/option"
)

func Ask(query string, path string) string {
	if os.Getenv("GEMINI_API_KEY") == "" {
		return "GEMINI_API_KEY is not set"
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Printf("Failed to create genai client: %v", err)
		return "Failed to create genai client"
	}
	defer client.Close()

	dsn := config.DSN{URL: path}
	s, err := datasource.Analyze(dsn)
	if err != nil {
		log.Printf("Failed to analyze schema: %v", err)
		return "Failed to analyze schema"
	}

	tpl, err := template.New("").Parse(templates.DefaultPromtTmpl)
	if err != nil {
		log.Printf("Failed to parse template: %v", err)
		return "Failed to parse template"
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, map[string]any{
		"DatabaseVersion": templates.DatabaseVersion(s),
		"QuoteStart":      "```sql",
		"QuoteEnd":        "```",
		"DDL":             templates.GenerateDDLRoughly(s),
		// "Question":        query,
	}); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return "Failed to execute template"
	}

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(buf.String()))
	if err != nil {
		log.Printf("Failed to generate content: %v", err)
		return "Failed to generate content"
	}
	answer := extractResponse(resp)
	return answer
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
