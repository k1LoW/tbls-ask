package schema

import (
	"fmt"
	"strings"

	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/datasource"
	"github.com/k1LoW/tbls/schema"
)

type Options struct {
	Includes   []string
	Excludes   []string
	Labels     []string
	Distance   int
}

func Load(strOrPath string, opts Options) (*schema.Schema, error) {
	var s *schema.Schema
	var err error

	if strings.HasPrefix(strOrPath, "{") || strings.HasPrefix(strOrPath, "/") {
		s, err = datasource.AnalyzeJSONStringOrFile(strOrPath)
	} else {
		dsn := config.DSN{URL: strOrPath}
		s, err = datasource.Analyze(dsn)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to analyze schema: %w", err)
	}

	if err := s.Filter(&schema.FilterOption{
		Include:       opts.Includes,
		Exclude:       opts.Excludes,
		IncludeLabels: opts.Labels,
		Distance:      opts.Distance,
	}); err != nil {
		return nil, fmt.Errorf("failed to filter schema: %w", err)
	}

	return s, nil
}
