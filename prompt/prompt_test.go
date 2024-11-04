package prompt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k1LoW/tbls/schema"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name    string
		schema  *schema.Schema
		want    string
		wantErr bool
	}{
		{
			name: "Generate prompt for simple schema",
			schema: &schema.Schema{
				Name: "testdb",
				Driver: &schema.Driver{
					Name:            "mysql",
					DatabaseVersion: "5.7",
				},
				Tables: []*schema.Table{
					{
						Name: "users",
						Columns: []*schema.Column{
							{Name: "id", Type: "int"},
							{Name: "name", Type: "varchar"},
						},
					},
				},
			},
			want: `Answer the questions in the Question assuming the following DDL.
Database is MySQL 5.7.

## DDL ( Data Definition Language )

` + "```sql" + `
CREATE TABLE users (
  id int NOT NULL,
  name varchar NOT NULL
);


` + "```" + `

`,
			wantErr: false,
		},
		{
			name: "Generate prompt for schema with multiple tables",
			schema: &schema.Schema{
				Name: "testdb",
				Driver: &schema.Driver{
					Name:            "mysql",
					DatabaseVersion: "5.7",
				},
				Tables: []*schema.Table{
					{
						Name: "users",
						Columns: []*schema.Column{
							{Name: "id", Type: "int"},
							{Name: "name", Type: "varchar"},
						},
					},
					{
						Name: "posts",
						Columns: []*schema.Column{
							{Name: "id", Type: "int"},
							{Name: "title", Type: "varchar"},
							{Name: "content", Type: "text"},
							{Name: "user_id", Type: "int"},
						},
					},
				},
			},
			want: `Answer the questions in the Question assuming the following DDL.
Database is MySQL 5.7.

## DDL ( Data Definition Language )

` + "```sql" + `
CREATE TABLE users (
  id int NOT NULL,
  name varchar NOT NULL
);

CREATE TABLE posts (
  id int NOT NULL,
  title varchar NOT NULL,
  content text NOT NULL,
  user_id int NOT NULL
);


` + "```" + `

`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.schema)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Generate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
