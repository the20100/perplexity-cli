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
  Set your API key via the PERPLEXITY_API_KEY environment variable,
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

func newClient() *client.Client {
	key := apiKey
	if key == "" {
		key = os.Getenv("PERPLEXITY_API_KEY")
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
