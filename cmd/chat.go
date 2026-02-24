package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vincentmaurin/perplexity-cli/client"
)

var chatCmd = &cobra.Command{
	Use:   "chat [message]",
	Short: "Chat with Perplexity AI (with web search grounding)",
	Long: `Send a message to Perplexity AI and get an AI-generated response grounded in web search.

If no message is given, starts an interactive multi-turn session.

Examples:
  perplexity chat "What is the capital of France?"
  perplexity chat --model sonar-pro "Explain quantum entanglement"
  perplexity chat --no-stream --json "Who won the 2024 election?"
  perplexity chat --system "You are a Go expert" "How do I use channels?"
  perplexity chat --recency week "Latest AI news"
  perplexity chat   # interactive mode`,
	RunE: runChat,
}

var (
	chatModel       string
	chatSystem      string
	chatStream      bool
	chatJSON        bool
	chatRecency     string
	chatDomains     []string
	chatMaxTokens   int
	chatTemperature float64
	chatCitations   bool
)

func init() {
	rootCmd.AddCommand(chatCmd)

	f := chatCmd.Flags()
	f.StringVarP(&chatModel, "model", "m", "sonar", "Model to use (e.g. sonar, sonar-pro)")
	f.StringVarP(&chatSystem, "system", "s", "", "System prompt")
	f.BoolVar(&chatStream, "stream", true, "Stream the response token by token")
	f.BoolVar(&chatJSON, "json", false, "Output raw JSON response (disables streaming)")
	f.StringVarP(&chatRecency, "recency", "r", "", "Search recency filter: hour, day, week, month, year")
	f.StringArrayVarP(&chatDomains, "domain", "d", nil, "Search domain filter (repeatable)")
	f.IntVar(&chatMaxTokens, "max-tokens", 0, "Maximum tokens in the response")
	f.Float64VarP(&chatTemperature, "temperature", "t", 0, "Sampling temperature (0.0-2.0, 0 = default)")
	f.BoolVar(&chatCitations, "citations", true, "Show citations at the end")
}

func runChat(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return runInteractiveChat(cmd)
	}
	return runSingleChat(cmd, strings.Join(args, " "))
}

func buildMessages(system, userMsg string) []client.Message {
	msgs := []client.Message{}
	if system != "" {
		msgs = append(msgs, client.Message{Role: "system", Content: system})
	}
	msgs = append(msgs, client.Message{Role: "user", Content: userMsg})
	return msgs
}

func buildChatRequest(msgs []client.Message) *client.ChatRequest {
	req := &client.ChatRequest{
		Model:    chatModel,
		Messages: msgs,
	}
	if chatRecency != "" {
		req.SearchRecencyFilter = chatRecency
	}
	if len(chatDomains) > 0 {
		req.SearchDomainFilter = chatDomains
	}
	if chatMaxTokens > 0 {
		req.MaxTokens = chatMaxTokens
	}
	if chatTemperature > 0 {
		req.Temperature = chatTemperature
	}
	return req
}

func printCitations(citations []client.Citation) {
	if len(citations) == 0 {
		return
	}
	fmt.Println("\nSources:")
	for i, c := range citations {
		if c.Title != "" {
			fmt.Printf("  [%d] %s\n      %s\n", i+1, c.Title, c.URL)
		} else {
			fmt.Printf("  [%d] %s\n", i+1, c.URL)
		}
	}
}

func runSingleChat(_ *cobra.Command, userMsg string) error {
	c := newClient()
	msgs := buildMessages(chatSystem, userMsg)
	req := buildChatRequest(msgs)

	if chatJSON {
		req.Stream = false
		resp, err := c.ChatComplete(req)
		if err != nil {
			return err
		}
		return printJSON(resp)
	}

	if chatStream {
		lastChunk, err := c.ChatStream(req, func(delta string) {
			fmt.Print(delta)
		})
		fmt.Println()
		if err != nil {
			return err
		}
		if chatCitations && lastChunk != nil {
			printCitations(lastChunk.Citations)
		}
		return nil
	}

	resp, err := c.ChatComplete(req)
	if err != nil {
		return err
	}
	if len(resp.Choices) > 0 {
		fmt.Println(resp.Choices[0].Message.Content)
	}
	if chatCitations {
		printCitations(resp.Citations)
	}
	return nil
}

func runInteractiveChat(_ *cobra.Command) error {
	c := newClient()

	fmt.Printf("Perplexity AI Chat (model: %s)\n", chatModel)
	fmt.Println("Type your message and press Enter. Use Ctrl+C or 'exit' to quit.\n")

	var history []client.Message
	if chatSystem != "" {
		history = append(history, client.Message{Role: "system", Content: chatSystem})
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.ToLower(line) == "exit" || strings.ToLower(line) == "quit" {
			break
		}

		history = append(history, client.Message{Role: "user", Content: line})
		req := buildChatRequest(history)

		fmt.Print("Assistant: ")
		var fullContent strings.Builder

		lastChunk, err := c.ChatStream(req, func(delta string) {
			fmt.Print(delta)
			fullContent.WriteString(delta)
		})
		fmt.Println()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			// Remove last user message so conversation stays consistent
			history = history[:len(history)-1]
			continue
		}

		if chatCitations && lastChunk != nil {
			printCitations(lastChunk.Citations)
		}

		history = append(history, client.Message{
			Role:    "assistant",
			Content: fullContent.String(),
		})

		fmt.Println()
	}

	fmt.Println("Goodbye!")
	return nil
}
