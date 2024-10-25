package client

import (
	"context"
	"testing"
	"github.com/google/go-cmp/cmp"
)

type StubLLMAgent struct{}

func (s *StubLLMAgent) ChatCompletionRequest(ctx context.Context, p string) (string, error) {
	return `stub response
` + "```sql" + `
SELECT * FROM users;
` + "```" + `
`, nil
}

func TestAsk(t *testing.T) {
	tests := []struct{
		name string
		prompt string
		querymode bool
		want string
		wantErr bool
	}{
		{
			name: "normal mode",
			prompt: "Hello, how are you?",
			querymode: false,
			want: `stub response
` + "```sql" + `
SELECT * FROM users;
` + "```" + `
`,
			wantErr: false,
		},
		{
			name: "query mode",
			prompt: "select users",
			querymode: true,
			want: "SELECT * FROM users;",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Agent: &StubLLMAgent{},
				Querymode: tt.querymode,
			}
			got, err := c.Ask(context.Background(), tt.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Ask() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
