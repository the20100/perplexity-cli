package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var contentCmd = &cobra.Command{
	Use:   "content <url> [url...]",
	Short: "Extract content from one or more URLs",
	Long: `Extract and retrieve full text content from web pages using the Perplexity AI Content API.

Examples:
  perplexity content https://example.com/article
  perplexity content https://site1.com/page https://site2.com/page
  perplexity content --json https://example.com`,
	Args: cobra.MinimumNArgs(1),
	RunE: runContent,
}

var contentJSON bool

func init() {
	rootCmd.AddCommand(contentCmd)
	contentCmd.Flags().BoolVar(&contentJSON, "json", false, "Output raw JSON response")
}

func runContent(cmd *cobra.Command, args []string) error {
	c := newClient()

	resp, err := c.GetContent(args)
	if err != nil {
		return err
	}

	if contentJSON {
		return printJSON(resp)
	}

	fmt.Printf("Extracted content from %d URL(s)\n\n", len(resp.Results))
	for i, r := range resp.Results {
		fmt.Printf("=== [%d] %s ===\n", i+1, r.URL)
		if r.Title != "" {
			fmt.Printf("Title: %s\n", r.Title)
		}
		fmt.Println()
		fmt.Println(strings.TrimSpace(r.Content))
		fmt.Println()
	}

	return nil
}
