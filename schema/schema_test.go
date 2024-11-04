package schema

import (
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	tblsschema "github.com/k1LoW/tbls/schema"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		strOrPath string
		opts      Options
		want      *tblsschema.Schema
		wantErr   bool
	}{
		{
			name:      "load JSON string",
			strOrPath: `{"name": "test", "tables": [{"name": "a", "comment": "table a", "columns": [{"name": "id", "type": "int"}]},{"name": "b", "comment": "table b", "columns": [{"name": "title", "type": "varchar"}]}]}`,
			opts:      Options{},
			want: &tblsschema.Schema{
				Name: "test",
				Tables: []*tblsschema.Table{
					{
						Name:    "a",
						Comment: "table a",
						Columns: []*tblsschema.Column{
							{
								Name: "id",
								Type: "int",
							},
						},
					},
					{
						Name:    "b",
						Comment: "table b",
						Columns: []*tblsschema.Column{
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
		{
			name:      "load with includes",
			strOrPath: `{"name": "test", "tables": [{"name": "a", "comment": "table a", "columns": [{"name": "id", "type": "int"}]},{"name": "b", "comment": "table b", "columns": [{"name": "title", "type": "varchar"}]}]}`,
			opts: Options{
				Includes: []string{"a"},
			},
			want: &tblsschema.Schema{
				Name: "test",
				Tables: []*tblsschema.Table{
					{
						Name:    "a",
						Comment: "table a",
						Columns: []*tblsschema.Column{
							{
								Name:            "id",
								Type:            "int",
								ParentRelations: []*tblsschema.Relation{},
								ChildRelations:  []*tblsschema.Relation{},
								// Add other fields with their default values
								Nullable:  false,
								PK:        false,
								FK:        false,
								Percents:  sql.NullFloat64{},
								Labels:    nil,
								HideForER: false,
							},
						},
					},
				},
				Relations: []*tblsschema.Relation{},
			},
			wantErr: false,
		},
		{
			name:      "load with excludes",
			strOrPath: `{"name": "test", "tables": [{"name": "a", "comment": "table a", "columns": [{"name": "id", "type": "int"}]},{"name": "b", "comment": "table b", "columns": [{"name": "title", "type": "varchar"}]}]}`,
			opts: Options{
				Excludes: []string{"b"},
			},
			want: &tblsschema.Schema{
				Name: "test",
				Tables: []*tblsschema.Table{
					{
						Name:    "a",
						Comment: "table a",
						Columns: []*tblsschema.Column{
							{
								Name:            "id",
								Type:            "int",
								ParentRelations: []*tblsschema.Relation{},
								ChildRelations:  []*tblsschema.Relation{},
								// Add other fields with their default values
								Nullable:  false,
								PK:        false,
								FK:        false,
								Percents:  sql.NullFloat64{},
								Labels:    nil,
								HideForER: false,
							},
						},
					},
				},
				Relations: []*tblsschema.Relation{},
			},
			wantErr: false,
		},
		{
			name:      "invalid JSON",
			strOrPath: `{"invalid": "json"`,
			opts:      Options{},
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "invalid DSN",
			strOrPath: "invalid:dsn",
			opts:      Options{},
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.strOrPath, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Load() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
