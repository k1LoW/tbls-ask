package templates

import (
	"fmt"
	"strings"

	"github.com/k1LoW/tbls/schema"
)

const (
	DefaultPromptTmpl = `Answer the questions in the Question assuming the following DDL.
{{ .DatabaseVersion }}

## DDL ( Data Definition Language )

{{ .QuoteStart }}
{{ .DDL }}
{{ .QuoteEnd }}
{{ if .ViewPoints }}

## ViewPoints (Sets of tables based on specific concerns)

{{ .ViewPoints }}

{{ end }}
`
)

func GenerateDDLRoughly(s *schema.Schema) string {
	var ddl string
	for _, t := range s.Tables {
		if t.Type == "VIEW" {
			continue
		}
		ddl += fmt.Sprintf("CREATE TABLE %s (\n", t.Name)
		td := []string{}
		for _, c := range t.Columns {
			d := fmt.Sprintf("  %s %s", c.Name, c.Type)
			if c.Default.String != "" {
				d += fmt.Sprintf(" DEFAULT %s", c.Default.String)
			}
			if !c.Nullable {
				d += " NOT NULL"
			}
			if c.Comment != "" {
				d += fmt.Sprintf(" COMMENT %q", c.Comment)
			}
			td = append(td, d)
		}
		for _, i := range t.Indexes {
			d := fmt.Sprintf("  %s", i.Def)
			td = append(td, d)
		}
		for _, c := range t.Constraints {
			switch c.Type {
			case "PRIMARY KEY", "UNIQUE KEY":
				continue
			default:
				d := fmt.Sprintf("  CONSTRAINT %s", c.Def)
				td = append(td, d)
			}
		}
		ddl += fmt.Sprintf("%s\n", strings.Join(td, ",\n"))
		if t.Comment != "" {
			ddl += fmt.Sprintf(") COMMENT = %q;\n\n", t.Comment)
		} else {
			ddl += ");\n\n"
		}
	}
	return ddl
}

func GenerateViewPoints(s *schema.Schema) string {
	var output strings.Builder
	for _, v := range s.Viewpoints {
		output.WriteString(fmt.Sprintf("Viewpoint: %s\n", v.Name))
		if v.Desc != "" {
			output.WriteString(fmt.Sprintf("- Description: %s\n", v.Desc))
		}
		if len(v.Labels) > 0 {
			output.WriteString(fmt.Sprintf("- Labels: %s\n", strings.Join(v.Labels, ", ")))
		}
		if len(v.Tables) > 0 {
			output.WriteString(fmt.Sprintf("- Tables: %s\n", strings.Join(v.Tables, ", ")))
		}
		if v.Distance > 0 {
			output.WriteString(fmt.Sprintf("- Distance: %d\n", v.Distance))
		}
		for _, g := range v.Groups {
			output.WriteString(fmt.Sprintf("- Group: %s\n", g.Name))
			if g.Desc != "" {
				output.WriteString(fmt.Sprintf("  - Description: %s\n", g.Desc))
			}
			if len(g.Labels) > 0 {
				output.WriteString(fmt.Sprintf("  - Labels: %s\n", strings.Join(g.Labels, ", ")))
			}
			if len(g.Tables) > 0 {
				output.WriteString(fmt.Sprintf("  - Tables: %s\n", strings.Join(g.Tables, ", ")))
			}
			if g.Color != "" {
				output.WriteString(fmt.Sprintf("  - Color: %s\n", g.Color))
			}
		}
	}
	return output.String()
}

func DatabaseVersion(s *schema.Schema) string {
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
