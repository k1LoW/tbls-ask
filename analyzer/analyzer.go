package analyzer

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/k1LoW/tbls-ask/client"
	"github.com/k1LoW/tbls-ask/templates"
	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/datasource"
	"github.com/k1LoW/tbls/schema"
)

type Analyzer struct {
	Schema *schema.Schema
}

func (a *Analyzer) AnalyzeSchema(strOrPath string) error {
	var s *schema.Schema
	var err error

	if strings.HasPrefix(strOrPath, "{") || strings.HasPrefix(strOrPath, "/") {
		s, err = datasource.AnalyzeJSONStringOrFile(strOrPath)
	} else {
		dsn := config.DSN{URL: strOrPath}
		s, err = datasource.Analyze(dsn)
	}
	if err != nil {
		return fmt.Errorf("failed to analyze schema: %w", err)
	}

	a.Schema = s

	return nil
}

func (a *Analyzer) FilterSchema(includes []string, excludes []string, labels []string, distance int) error {
	if err := a.Schema.Filter(&schema.FilterOption{
		Include:       includes,
		Exclude:       excludes,
		IncludeLabels: labels,
		Distance:      distance,
	}); err != nil {
		return fmt.Errorf("failed to filter schema: %w", err)
	}
	return nil
}

func (a *Analyzer) GeneratePrompt(q string, querymode bool) (string, error) {
	var promptTmpl string
	if querymode {
		promptTmpl = templates.DefaultQueryPromptTmpl
	} else {
		promptTmpl = templates.DefaultPromptTmpl
	}
	tpl, err := template.New("").Parse(promptTmpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, map[string]any{
		"DatabaseVersion": templates.DatabaseVersion(a.Schema),
		"QuoteStart":      "```sql",
		"QuoteEnd":        "```",
		"DDL":             templates.GenerateDDLRoughly(a.Schema),
		"Viewpoints":      templates.GenerateViewPoints(a.Schema),
		"Question":        q,
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a *Analyzer) ExtractRelevantTables(ctx context.Context, c *client.Client, query string) ([]string, error) {
	var info string
	tableNames := make([]string, 0, len(a.Schema.Tables))
	for _, t := range a.Schema.Tables {
		if t.Type == "VIEW" || t.Type == "MATERIALIZED VIEW" {
			continue
		}
		tableNames = append(tableNames, t.Name)
		if t.Comment != "" {
			info += fmt.Sprintf("%s: %s\n", t.Name, t.Comment)
		} else {
			info += fmt.Sprintf("%s\n", t.Name)
		}
	}

	prompt := fmt.Sprintf(templates.RelevantTablesPromptTmpl, info, query)
	resp, err := c.Agent.ChatCompletionRequest(ctx, prompt)
	if err != nil {
		return nil, err
	}
	if resp == "" {
		return nil, fmt.Errorf("failed to extract relevant tables: empty response")
	}

	relevantTables := strings.Split(resp, ",")
	for i, t := range relevantTables {
		relevantTables[i] = strings.TrimSpace(t)
		found := false
		for _, tableName := range tableNames {
			if tableName == relevantTables[i] {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("failed to extract relevant tables: %s not found", relevantTables[i])
		}
	}

	return relevantTables, nil
}
