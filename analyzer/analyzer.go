package analyzer

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/k1LoW/tbls-ask/templates"
	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/datasource"
	"github.com/k1LoW/tbls/schema"
)

type Analyzer struct {
	Schema *schema.Schema
}

func (a *Analyzer) AnalyzeSchema(strOrPath string, includes []string, excludes []string, labels []string) error {
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

	if err := s.Filter(&schema.FilterOption{
		Include:       includes,
		Exclude:       excludes,
		IncludeLabels: labels,
	}); err != nil {
		return fmt.Errorf("failed to filter schema: %w", err)
	}

	a.Schema = s

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
