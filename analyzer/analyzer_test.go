package analyzer

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k1LoW/tbls-ask/client"
	"github.com/k1LoW/tbls/schema"
)

func TestAnalayzeSchema(t *testing.T) {
	tests := []struct {
		name      string
		strOrPath string
		want      *schema.Schema
		wantErr   bool
	}{
		{
			name:      "analyze string",
			strOrPath: `{"name": "test", "tables": [{"name": "a", "comment": "table a", "columns": [{"name": "id", "type": "int"}]},{"name": "b", "comment": "table b", "columns": [{"name": "title", "type": "varchar"}]}]}`,
			want: &schema.Schema{
				Name: "test",
				Tables: []*schema.Table{
					{
						Name:    "a",
						Comment: "table a",
						Columns: []*schema.Column{
							{
								Name: "id",
								Type: "int",
							},
						},
					},
					{
						Name:    "b",
						Comment: "table b",
						Columns: []*schema.Column{
							{
								Name: "title",
								Type: "varchar",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analyzer{}
			err := a.AnalyzeSchema(tt.strOrPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, a.Schema); diff != "" {
				t.Errorf("AnalyzeSchema() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFilterSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   *schema.Schema
		includes []string
		excludes []string
		labels   []string
		distance int
		want     *schema.Schema
		wantErr  bool
	}{
		{
			name: "filter tables",
			schema: &schema.Schema{
				Name: "test",
				Tables: []*schema.Table{
					{Name: "users"},
					{Name: "posts"},
					{Name: "comments"},
				},
				Relations: []*schema.Relation{},
			},
			includes: []string{"users", "posts"},
			excludes: []string{},
			labels:   []string{},
			distance: 0,
			want: &schema.Schema{
				Name: "test",
				Tables: []*schema.Table{
					{Name: "users"},
					{Name: "posts"},
				},
				Relations: []*schema.Relation{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analyzer{Schema: tt.schema}
			err := a.FilterSchema(tt.includes, tt.excludes, tt.labels, tt.distance)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, a.Schema); diff != "" {
				t.Errorf("FilterSchema() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGeneratePrompt(t *testing.T) {
	tests := []struct {
		name      string
		q         string
		querymode bool
		want      string
		wantErr   bool
	}{
		{
			name:      "normalmode",
			q:         "select users",
			querymode: false,
			want: `Answer the questions in the Question assuming the following DDL.
Database is MySQL 5.7.

## DDL ( Data Definition Language )

` + "```sql" + `
CREATE TABLE users (
  id int NOT NULL,
  name varchar NOT NULL
) COMMENT = "table users";


` + "```" + `

## Question
select users
`,
			wantErr: false,
		},
		{
			name:      "querymode",
			q:         "select users",
			querymode: true,
			want: `Answer the SQL query in the "Explanation of the query to be created" section, assuming the database was created with the following DDL.
Database is MySQL 5.7.

## DDL ( Data Definition Language )

` + "```sql" + `
CREATE TABLE users (
  id int NOT NULL,
  name varchar NOT NULL
) COMMENT = "table users";


` + "```" + `

## Explanation of the query to be created
select users
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analyzer{
				Schema: &schema.Schema{
					Name: "test",
					Driver: &schema.Driver{
						Name:            "mysql",
						DatabaseVersion: "5.7",
					},
					Tables: []*schema.Table{
						{
							Name:    "users",
							Comment: "table users",
							Columns: []*schema.Column{
								{
									Name: "id",
									Type: "int",
								},
								{
									Name: "name",
									Type: "varchar",
								},
							},
						},
					},
				},
			}
			got, err := a.GeneratePrompt(tt.q, tt.querymode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GeneratePrompt() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestExtractRelevantTables(t *testing.T) {
	tests := []struct {
		name    string
		schema  *schema.Schema
		query   string
		want    []string
		wantErr bool
	}{
		{
			name: "extract relevant tables",
			schema: &schema.Schema{
				Tables: []*schema.Table{
					{Name: "users", Type: "BASE TABLE"},
					{Name: "posts", Type: "BASE TABLE"},
					{Name: "comments", Type: "BASE TABLE"},
					{Name: "user_view", Type: "VIEW"},
				},
			},
			query:   "Select all users and their posts",
			want:    []string{"users", "posts"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analyzer{Schema: tt.schema}
			mockClient := &client.Client{
				Agent: &MockAgent{
					Response: "users, posts",
				},
			}
			got, err := a.ExtractRelevantTables(context.Background(), mockClient, tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractRelevantTables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ExtractRelevantTables() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type MockAgent struct {
	Response string
}

func (m *MockAgent) ChatCompletionRequest(ctx context.Context, prompt string) (string, error) {
	return m.Response, nil
}
