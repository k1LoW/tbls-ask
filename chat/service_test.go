package chat

import (
	"context"
	"testing"
)

func TestNewService(t *testing.T) {
	t.Run("with OPENAI_API_KEY", func(t *testing.T) {
		t.Setenv("OPENAI_API_KEY", "test-openai-key")

		tests := []struct {
			name  string
			model string
		}{
			{"GPT model alias", "chat-latest"},
			{"O-series model", "o1-mini"},
			{"Claude model via OpenRouter", "anthropic/claude-sonnet-latest"},
			{"Gemini model via compatibility endpoint", "gemini-flash-latest"},
			{"Custom local model", "my-custom-model"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := NewService(tt.model)
				if err != nil {
					t.Errorf("NewService() error = %v, want nil", err)
				}
			})
		}
	})

	t.Run("with OPENAI_BASE_URL", func(t *testing.T) {
		t.Setenv("OPENAI_API_KEY", "test-openai-key")
		t.Setenv("OPENAI_BASE_URL", "https://openrouter.ai/api/v1")

		_, err := NewService("anthropic/claude-sonnet-latest")
		if err != nil {
			t.Errorf("NewService() with OPENAI_BASE_URL error = %v, want nil", err)
		}
	})

	t.Run("without API key", func(t *testing.T) {
		t.Setenv("OPENAI_API_KEY", "")
		t.Setenv("AZURE_OPENAI_ENDPOINT", "")

		_, err := NewService("chat-latest")
		if err == nil {
			t.Errorf("NewService() want error when OPENAI_API_KEY is not set, got nil")
		}
	})

	t.Run("with Azure OpenAI", func(t *testing.T) {
		t.Setenv("AZURE_OPENAI_ENDPOINT", "https://test.openai.azure.com")
		t.Setenv("AZURE_OPENAI_KEY", "test-azure-key")

		_, err := NewService("my-deployment")
		if err != nil {
			t.Errorf("NewService() with Azure error = %v, want nil", err)
		}
	})
}

func TestService_Ask(t *testing.T) {
	mockClient := &MockClient{}
	service := &Service{client: mockClient}

	tests := []struct {
		name      string
		messages  []Message
		queryMode bool
		want      string
		wantErr   bool
	}{
		{
			name: "Normal response",
			messages: []Message{
				{Role: "user", Content: "Hello"},
			},
			queryMode: false,
			want:      "Hello, how can I help you?",
			wantErr:   false,
		},
		{
			name: "Query mode",
			messages: []Message{
				{Role: "user", Content: "Show me all users"},
			},
			queryMode: true,
			want:      "SELECT * FROM users;",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.AskFunc = func(ctx context.Context, messages []Message) (string, error) {
				if tt.queryMode {
					return "Here's the query:\n```sql\nSELECT * FROM users;\n```\n", nil
				}
				return tt.want, nil
			}

			got, err := service.Ask(context.Background(), tt.messages, tt.queryMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Ask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.Ask() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockClient struct {
	AskFunc func(ctx context.Context, messages []Message) (string, error)
}

func (m *MockClient) Ask(ctx context.Context, messages []Message) (string, error) {
	return m.AskFunc(ctx, messages)
}
