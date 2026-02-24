# perplexity-cli

A Go CLI for the [Perplexity AI API](https://docs.perplexity.ai) — web search, AI chat with grounding, and URL content extraction, all from your terminal.

## Install

```bash
git clone https://github.com/vincentmaurin/perplexity-cli
cd perplexity-cli
go build -o perplexity .
# optionally move to PATH
mv perplexity /usr/local/bin/
```

## Auth

Get your API key from [perplexity.ai/settings/api](https://www.perplexity.ai/settings/api) and export it:

```bash
export PERPLEXITY_API_KEY=pplx-...
```

Or pass it inline on any command:

```bash
perplexity --api-key pplx-... search "query"
```

---

## Commands

### `search` — Real-time web search

```bash
perplexity search <query> [flags]
```

**Flags**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--results <n>` | `-n` | `10` | Number of results (1–20) |
| `--mode <mode>` | `-m` | | Search mode: `web`, `academic`, `sec` |
| `--recency <period>` | `-r` | | Recency filter: `hour` `day` `week` `month` `year` |
| `--country <code>` | `-c` | | ISO 3166-1 alpha-2 country code (e.g. `us`, `fr`) |
| `--domain <domain>` | `-d` | | Domain allowlist, repeatable |
| `--lang <code>` | `-l` | | Language filter (ISO 639-1), repeatable |
| `--after <date>` | | | Results after this date (`YYYY-MM-DD`) |
| `--before <date>` | | | Results before this date (`YYYY-MM-DD`) |
| `--max-tokens <n>` | | | Max tokens across all results |
| `--snippet` | | | Compact output: URL + snippet only |
| `--json` | | | Raw JSON output |

**Examples**

```bash
# Basic
perplexity search "latest Go releases"

# Academic, top 5
perplexity search --results 5 --mode academic "quantum computing"

# News from this week, US sources
perplexity search --recency week --country us "AI news"

# Trusted scientific domains only
perplexity search --domain nature.com --domain science.org "climate research"

# Date range
perplexity search --after 2024-01-01 --before 2024-12-31 "election results"

# Compact one-liner output
perplexity search --snippet "Go generics"

# Pipe-friendly JSON
perplexity search --json "Go 1.22 features" | jq '.results[].url'
```

---

### `chat` — AI chat with web grounding

```bash
perplexity chat [message] [flags]
```

Omit the message to enter an **interactive multi-turn session**.

Responses are **streamed by default** — you see tokens as they arrive.

**Flags**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--model <model>` | `-m` | `sonar` | Model: `sonar`, `sonar-pro` |
| `--system <prompt>` | `-s` | | System prompt |
| `--stream` | | `true` | Stream tokens as they arrive |
| `--recency <period>` | `-r` | | Recency filter on web search |
| `--domain <domain>` | `-d` | | Domain filter (repeatable) |
| `--max-tokens <n>` | | | Max tokens in response |
| `--temperature <f>` | `-t` | | Sampling temperature (0.0–2.0) |
| `--citations` | | `true` | Print source citations after the answer |
| `--json` | | | Raw JSON (disables streaming) |

**Examples**

```bash
# Simple question
perplexity chat "What is the capital of France?"

# Better model
perplexity chat --model sonar-pro "Explain quantum entanglement simply"

# System prompt for a persona
perplexity chat --system "You are a senior Go engineer" "Best practices for error handling?"

# Only draw from recent sources
perplexity chat --recency week "Latest LLM releases"

# Restrict sources to specific domains
perplexity chat --domain arxiv.org "Recent advances in diffusion models"

# No streaming, get full JSON
perplexity chat --no-stream --json "Summarise the history of the internet"

# Interactive multi-turn conversation
perplexity chat
```

**Interactive mode** — type `exit` or press `Ctrl+C` to quit. The full conversation history is kept so follow-up questions work naturally.

---

### `content` — Extract text from URLs

```bash
perplexity content <url> [url...] [flags]
```

**Flags**

| Flag | Description |
|------|-------------|
| `--json` | Raw JSON output |

**Examples**

```bash
# Single article
perplexity content https://go.dev/blog/go1.22

# Multiple pages at once
perplexity content https://site1.com/page https://site2.com/page

# JSON for scripting
perplexity content --json https://example.com | jq '.results[0].content'
```

---

## Models

| Model | Description |
|-------|-------------|
| `sonar` | Fast, cost-efficient, real-time web search |
| `sonar-pro` | Stronger reasoning, longer context, better citations |

---

## Project structure

```
perplexity-cli/
├── main.go              # Entry point
├── cmd/
│   ├── root.go          # Root command, shared helpers (newClient, printJSON)
│   ├── search.go        # search command
│   ├── chat.go          # chat command — streaming + interactive mode
│   └── content.go       # content command
└── client/
    ├── client.go        # HTTP client (Search, ChatComplete, ChatStream, GetContent)
    └── types.go         # All request / response types
```

## Tech

- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Viper](https://github.com/spf13/viper) — config management
- Native `net/http` — no extra HTTP client dependency
- Server-Sent Events parsed manually for streaming

## License

MIT
