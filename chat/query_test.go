package chat

import (
	"testing"
)

func TestExtractQuery(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid query",
			input:   "Here's the query:\n```sql\nSELECT * FROM users;\n```\n",
			want:    "SELECT * FROM users;",
			wantErr: false,
		},
		{
			name:    "No SQL block",
			input:   "There's no SQL query here.",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractQuery(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractQuery() = %q, want %q", got, tt.want)
			}
		})
	}
}
