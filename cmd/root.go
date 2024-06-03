/*
Copyright © 2023 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/k1LoW/tbls-ask/gemini"
	"github.com/k1LoW/tbls-ask/openai"
	"github.com/k1LoW/tbls-ask/version"
	"github.com/k1LoW/tbls/datasource"
	"github.com/k1LoW/tbls/schema"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var (
	query    bool
	model    string
	tables   []string
	includes []string
	excludes []string
	labels   []string
)

// OpenAI または Gemini を型として持つ変数を作成する
// この変数は、OpenAI または Gemini の Ask メソッドを呼び出すために使用される
var m interface {
	Ask(ctx context.Context, q string, s *schema.Schema) (string, error)
	AskQuery(ctx context.Context, q string, s *schema.Schema) (string, error)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "tbls-ask",
	Short:        "ask LLM using the datasource",
	Long:         `ask LLM using the datasource.`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if strings.HasPrefix(model, "gpt") {
			if os.Getenv("OPENAI_API_KEY") == "" {
				return errors.New("OPENAI_API_KEY is not set")
			}
		} else if strings.HasPrefix(model, "gemini") {
			if os.Getenv("GEMINI_API_KEY") == "" {
				return errors.New("GEMINI_API_KEY is not set")
			}
		} else {
			return errors.New("model is not supported")
		}
		if os.Getenv("TBLS_SCHEMA") == "" {
			return errors.New("TBLS_SCHEMA is not set")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		// model が gpt から始まる場合は OpenAI を使う
		if strings.HasPrefix(model, "gpt") {
			m = openai.New(os.Getenv("OPENAI_API_KEY"), model)
		} else if strings.HasPrefix(model, "gemini") {
			m = gemini.New(os.Getenv("GEMINI_API_KEY"), model)
		}
		q := strings.Join(args, " ")
		s, err := datasource.AnalyzeJSONStringOrFile(os.Getenv("TBLS_SCHEMA"))
		if err != nil {
			return fmt.Errorf("failed to analyze schema: %w", err)
		}
		includes = lo.Uniq(append(includes, tables...))
		if err := s.Filter(&schema.FilterOption{
			Include:       includes,
			Exclude:       excludes,
			IncludeLabels: labels,
		}); err != nil {
			return fmt.Errorf("failed to filter schema: %w", err)
		}

		var a string
		if query {
			a, err = m.AskQuery(ctx, q, s)
			if err != nil {
				return err
			}
		} else {
			a, err = m.Ask(ctx, q, s)
			if err != nil {
				return err
			}
		}
		cmd.Println(a)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	log.SetOutput(io.Discard)
	if env := os.Getenv("DEBUG"); env != "" {
		debug, err := os.Create(fmt.Sprintf("%s.debug", version.Name))
		if err != nil {
			rootCmd.PrintErr(err)
			os.Exit(1)
		}
		log.SetOutput(debug)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringSliceVarP(&tables, "table", "", []string{}, "target table (tables to include)")
	rootCmd.Flags().StringSliceVarP(&includes, "include", "", []string{}, "tables to include")
	rootCmd.Flags().StringSliceVarP(&excludes, "exclude", "", []string{}, "tables to exclude")
	rootCmd.Flags().StringSliceVarP(&labels, "label", "", []string{}, "table labels to be included")
	rootCmd.Flags().BoolVarP(&query, "query", "q", false, "ask OpenAI for query using the datasource")
	rootCmd.Flags().StringVarP(&model, "model", "m", openai.DefaultModelChat, "model to be used")
}
