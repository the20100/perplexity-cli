package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/the20100/perplexity-cli/client"
)

var apiKey string

var rootCmd = &cobra.Command{
	Use:   "perplexity",
	Short: "Perplexity AI CLI",
	Long: `A command-line interface for the Perplexity AI API.

Supports web search, AI chat with web grounding, and URL content extraction.

Configuration:
  Set your API key via the PERPLEXITY_API_KEY environment variable (or aliases:
  PERPLEXITY_KEY, PERPLEXITY_API, API_KEY_PERPLEXITY, ...),
  or pass it with --api-key.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Perplexity API key (or set PERPLEXITY_API_KEY)")
}

// resolveEnv returns the value of the first non-empty environment variable from the given names.
func resolveEnv(names ...string) string {
	for _, name := range names {
		if v := os.Getenv(name); v != "" {
			return v
		}
	}
	return ""
}

func newClient() *client.Client {
	key := apiKey
	if key == "" {
		key = resolveEnv(
			"PERPLEXITY_API_KEY", "PERPLEXITY_KEY", "PERPLEXITY_API", "API_KEY_PERPLEXITY", "API_PERPLEXITY", "PERPLEXITY_PK", "PERPLEXITY_PUBLIC",
			"PERPLEXITY_API_SECRET", "PERPLEXITY_SECRET_KEY", "PERPLEXITY_API_SECRET_KEY", "PERPLEXITY_SECRET", "SECRET_PERPLEXITY", "API_SECRET_PERPLEXITY", "SK_PERPLEXITY", "PERPLEXITY_SK",
		)
	}
	if key == "" {
		fmt.Fprintln(os.Stderr, "error: API key required. Set PERPLEXITY_API_KEY or use --api-key")
		os.Exit(1)
	}
	return client.New(key)
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
