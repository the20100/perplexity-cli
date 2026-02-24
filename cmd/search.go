package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vincentmaurin/perplexity-cli/client"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search the web using Perplexity AI",
	Long: `Perform a real-time web search using the Perplexity AI Search API.

Examples:
  perplexity search "latest Go releases"
  perplexity search --results 5 --mode academic "quantum computing"
  perplexity search --recency week --country us "AI news"
  perplexity search --domain nature.com --domain science.org "climate research"
  perplexity search --after 2024-01-01 --before 2024-12-31 "election results"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSearch,
}

var (
	searchResults      int
	searchMode         string
	searchRecency      string
	searchCountry      string
	searchDomains      []string
	searchLangs        []string
	searchAfter        string
	searchBefore       string
	searchMaxTokens    int
	searchJSON         bool
	searchSnippetOnly  bool
)

func init() {
	rootCmd.AddCommand(searchCmd)

	f := searchCmd.Flags()
	f.IntVarP(&searchResults, "results", "n", 10, "Number of results (1-20)")
	f.StringVarP(&searchMode, "mode", "m", "", "Search mode: web, academic, sec")
	f.StringVarP(&searchRecency, "recency", "r", "", "Recency filter: hour, day, week, month, year")
	f.StringVarP(&searchCountry, "country", "c", "", "Country code (ISO 3166-1 alpha-2, e.g. us, fr)")
	f.StringArrayVarP(&searchDomains, "domain", "d", nil, "Domain filter (repeatable, e.g. --domain nature.com)")
	f.StringArrayVarP(&searchLangs, "lang", "l", nil, "Language filter (repeatable, e.g. --lang en)")
	f.StringVar(&searchAfter, "after", "", "Only results after date (YYYY-MM-DD)")
	f.StringVar(&searchBefore, "before", "", "Only results before date (YYYY-MM-DD)")
	f.IntVar(&searchMaxTokens, "max-tokens", 0, "Max tokens across all results")
	f.BoolVar(&searchJSON, "json", false, "Output raw JSON response")
	f.BoolVar(&searchSnippetOnly, "snippet", false, "Show only URL and snippet (compact)")
}

func runSearch(cmd *cobra.Command, args []string) error {
	c := newClient()

	query := strings.Join(args, " ")

	req := &client.SearchRequest{
		Query:      query,
		MaxResults: searchResults,
	}
	if searchMode != "" {
		req.SearchMode = searchMode
	}
	if searchRecency != "" {
		req.SearchRecencyFilter = searchRecency
	}
	if searchCountry != "" {
		req.Country = searchCountry
	}
	if len(searchDomains) > 0 {
		req.SearchDomainFilter = searchDomains
	}
	if len(searchLangs) > 0 {
		req.SearchLanguageFilter = searchLangs
	}
	if searchAfter != "" {
		req.SearchAfterDateFilter = searchAfter
	}
	if searchBefore != "" {
		req.SearchBeforeDateFilter = searchBefore
	}
	if searchMaxTokens > 0 {
		req.MaxTokens = searchMaxTokens
	}

	resp, err := c.Search(req)
	if err != nil {
		return err
	}

	if searchJSON {
		return printJSON(resp)
	}

	fmt.Printf("Search results for: %q\n", query)
	if len(resp.Results) == 0 {
		fmt.Println("No results found.")
		return nil
	}
	fmt.Printf("Found %d result(s)\n\n", len(resp.Results))

	for i, r := range resp.Results {
		if searchSnippetOnly {
			fmt.Printf("[%d] %s\n    %s\n\n", i+1, r.URL, r.Snippet)
			continue
		}

		fmt.Printf("[%d] %s\n", i+1, r.Title)
		fmt.Printf("    URL: %s\n", r.URL)
		if r.Date != "" {
			fmt.Printf("    Date: %s\n", r.Date)
		}
		fmt.Printf("    %s\n\n", r.Snippet)
	}

	return nil
}
