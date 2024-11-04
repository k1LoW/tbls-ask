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

	"github.com/k1LoW/tbls-ask/chat"
	"github.com/k1LoW/tbls-ask/prompt"
	"github.com/k1LoW/tbls-ask/schema"
	"github.com/k1LoW/tbls-ask/version"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var (
	queryMode bool
	model     string
	tables    []string
	includes  []string
	excludes  []string
	labels    []string
	distance  int
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
		s := os.Getenv("TBLS_SCHEMA") // this env var is to be set by tbls
		if s == "" {
			return fmt.Errorf("TBLS_SCHEMA is not set")
		}

		opts := schema.Options{
			Includes: lo.Uniq(append(includes, tables...)),
			Excludes: excludes,
			Labels:   labels,
			Distance: distance,
		}

		schema, err := schema.Load(s, opts)
		if err != nil {
			return err
		}

		prompt, err := prompt.Generate(schema)
		if err != nil {
			return err
		}

		service, err := chat.NewService(model)
		if err != nil {
			return err
		}

		messages := []chat.Message{
			{
				Role:    "system",
				Content: "You are a database expert. You are given a database schema and a question. Answer the question based on the schema.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
			{
				Role:    "user",
				Content: q,
			},
		}

		response, err := service.Ask(ctx, messages, queryMode)
		if err != nil {
			return err
		}
		cmd.Println(response)
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
	rootCmd.Flags().BoolVarP(&queryMode, "query", "q", false, "ask OpenAI for query using the datasource")
	rootCmd.Flags().StringVarP(&model, "model", "m", "gpt-4o", "model to be used")
	rootCmd.Flags().IntVarP(&distance, "distance", "d", 1, "distance between tables to be included")
}
