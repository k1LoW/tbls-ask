package chat

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/k1LoW/repin"
)

const (
	quoteStart = "```sql"
	quoteEnd   = "```"
)

func ExtractQuery(resp string) (string, error) {
	if !strings.Contains(resp, quoteStart) {
		return "", fmt.Errorf("failed to pick query from answer: %s", resp)
	}

	if !strings.HasSuffix(resp, "\n") {
		resp += "\n"
	}

	src := strings.NewReader(resp)
	query := new(bytes.Buffer)
	if _, err := repin.Pick(src, quoteStart, quoteEnd, true, query); err != nil {
		return "", fmt.Errorf("failed to pick query from answer: %w\nanswer: %s", err, resp)
	}

	return query.String(), nil
}
