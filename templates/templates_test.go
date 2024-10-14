package templates

import (
	"fmt"
	"testing"

	"github.com/k1LoW/tbls/schema"
)

func TestGenerateDDLRoughly(t *testing.T) {
	tests := []struct {
		s    *schema.Schema
		want string
	}{
		{
			&schema.Schema{
				Name: "test",
				Tables: []*schema.Table{
					{
						Name: "test",
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
			`CREATE TABLE test (
  id int NOT NULL,
  name varchar NOT NULL
);

`,
		},
		{
			&schema.Schema{
				Name: "test",
				Tables: []*schema.Table{
					{
						Name:    "test",
						Comment: "test table",
						Columns: []*schema.Column{
							{
								Name:    "id",
								Type:    "int",
								Comment: "ID",
							},
							{
								Name:     "name",
								Type:     "varchar",
								Nullable: true,
							},
						},
					},
				},
			},
			`CREATE TABLE test (
  id int NOT NULL COMMENT "ID",
  name varchar
) COMMENT = "test table";

`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := GenerateDDLRoughly(tt.s)
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestGenerateViewPoints(t *testing.T) {
	tests := []struct {
		s    *schema.Schema
		want string
	}{
		{
			&schema.Schema{
				Name: "test",
				Tables: []*schema.Table{
					{
						Name: "test",
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
				Viewpoints: []*schema.Viewpoint{
					{
						Name:        "test",
						Desc: "test viewpoint",
						Tables: []string{
							"test",
						},
					},
				},
			},
			`Viewpoint: test
- Description: test viewpoint
- Tables: test
`,
	},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := GenerateViewPoints(tt.s)
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}
