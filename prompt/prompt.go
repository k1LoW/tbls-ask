// prompt/prompt.go
package prompt

import (
	"bytes"
	"text/template"

	"github.com/k1LoW/tbls-ask/templates"
	"github.com/k1LoW/tbls/schema"
)

func Generate(s *schema.Schema) (string, error) {
	promptTmpl := templates.DefaultPromptTmpl
	tpl, err := template.New("").Parse(promptTmpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, map[string]any{
		"DatabaseVersion": templates.DatabaseVersion(s),
		"QuoteStart":      "```sql",
		"QuoteEnd":        "```",
		"DDL":             templates.GenerateDDLRoughly(s),
		"Viewpoints":      templates.GenerateViewPoints(s),
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
