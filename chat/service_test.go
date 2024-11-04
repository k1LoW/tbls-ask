package chat

import (
	"context"
	"testing"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name    string
		model   string
		wantErr bool
	}{
		{"GPT model", "gpt-3.5-turbo", false},
		{"Gemini model", "gemini-pro", false},
		{"Unsupported model", "unsupported-model", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewService(tt.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
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
