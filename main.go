package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Reads all data from stdin. Returns empty string if there's no piped input.
func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	// If input is from a pipe or file, ModeCharDevice will be false
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}
	var b strings.Builder
	r := bufio.NewReader(os.Stdin)
	for {
		chunk, err := r.ReadString('\n')
		b.WriteString(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(b.String()), nil
}

func main() {
	// Read query from stdin
	query, err := readStdin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read stdin:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(query) == "" {
		fmt.Fprintln(os.Stderr, "no input provided on stdin; pipe a prompt, e.g.: echo 'Hello' | go run main.go")
		os.Exit(2)
	}

	// Load configuration from environment
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	baseURL := strings.TrimSpace(os.Getenv("OPENAI_BASE_URL"))
	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "missing OPENAI_API_KEY environment variable")
		os.Exit(3)
	}
	if model == "" {
		fmt.Fprintln(os.Stderr, "missing OPENAI_MODEL environment variable")
		os.Exit(4)
	}

	// Build client options
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)
	ctx := context.Background()

	// Use the Responses API to perform a basic text response to the user's prompt
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(query),
		},
		Model: model,
		Seed:  openai.Int(0),
	}
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	// Print the model's output to stdout
	fmt.Print(completion.Choices[0].Message.Content)
}
