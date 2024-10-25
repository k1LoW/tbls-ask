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
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/k1LoW/tbls-ask/analyzer"
	"github.com/k1LoW/tbls-ask/client"
	"github.com/k1LoW/tbls-ask/internal/gemini"
	"github.com/k1LoW/tbls-ask/internal/openai"
	"github.com/k1LoW/tbls-ask/version"
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "tbls-ask",
	Short:        "ask database info to LLM using provided table document",
	Long:         `ask database information to LLM using provided table document.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		q := strings.Join(args, " ")

		includes = lo.Uniq(append(includes, tables...)) // tables and includes are eqivalent
		s := os.Getenv("TBLS_SCHEMA") // this env var is to be set by tbls
		if s == "" {
			return fmt.Errorf("TBLS_SCHEMA is not set")
		}
		var a analyzer.Analyzer
		err := a.AnalyzeSchema(s, includes, excludes, labels)
		if err != nil {
			return err
		}

		p, err := a.GeneratePrompt(q, query)
		if err != nil {
			return err
		}

		var agent client.LLMAgent
		if strings.HasPrefix(model, "gpt") {
			agent, err = openai.NewClient(model)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(model, "gemini") {
			agent, err = gemini.NewClient(ctx, model)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported model: %s", model)
		}

		c := client.Client{
			Agent:     agent,
			Querymode: query,
		}

		answer, err := c.Ask(ctx, p)
		if err != nil {
			return err
		}
		cmd.Println(answer)
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
	rootCmd.Flags().StringVarP(&model, "model", "m", "gpt-4o", "model to be used")
}
