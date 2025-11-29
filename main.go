package main

import (
	"agent-sandbox/adapters"
	"agent-sandbox/llm"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	query := adapters.ParseQuery()

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
			openai.SystemMessage(llm.SystemPrompt),
			openai.UserMessage(query),
		},
		Model: model,
		Seed:  openai.Int(0),
		Tools: []openai.ChatCompletionToolUnionParam{
			adapters.ReadFileTool,
			adapters.WriteFileTool,
		},
	}

	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "read_file" {
			log.Printf("reading file: ")
			var args map[string]interface{}
			err := json.Unmarshal([]byte(getToolCall(toolCall.Function.Arguments)), &args)

			if err != nil {
				log.Printf("error unmarshalling arguments: %v", err)
				params.Messages = append(params.Messages, openai.ToolMessage(err.Error(), toolCall.ID))
				continue
			}
			path := args["path"].(string)
			log.Println(path)

			fileContents, err := adapters.ReadFile(path)
			if err != nil {
				log.Printf("error opening file: %v", err)
				params.Messages = append(params.Messages, openai.ToolMessage(err.Error(), toolCall.ID))
				continue
			}

			params.Messages = append(params.Messages, openai.ToolMessage(fileContents, toolCall.ID))
		} else if toolCall.Function.Name == "write_file" {
			log.Printf("writing file: ")
			var args map[string]interface{}
			err := json.Unmarshal([]byte(getToolCall(toolCall.Function.Arguments)), &args)
			if err != nil {
				log.Printf("error unmarshalling arguments: %v", err)
				params.Messages = append(params.Messages, openai.ToolMessage(err.Error(), toolCall.ID))
				continue
			}
			path := args["path"].(string)
			data := args["data"].(string)
			log.Println(path)
			log.Println(data)

			err = adapters.WriteFile(path, data)
			if err != nil {
				log.Printf("error writing file: %v", err)
				params.Messages = append(params.Messages, openai.ToolMessage(err.Error(), toolCall.ID))
			}
		} else {
			message := "Unknown tool: " + toolCall.Function.Name
			log.Printf(message)
			params.Messages = append(params.Messages, openai.ToolMessage(message, toolCall.ID))
		}

		completion, err = client.Chat.Completions.New(ctx, params)
		if err != nil {
			panic(err)
		}
		toolCalls = completion.Choices[0].Message.ToolCalls
	}

	// Print the model's output to stdout
	log.Println("Got model output:")
	log.Println(completion.Choices[0].Message.Content)
	fmt.Print(completion.Choices[0].Message.Content)
}

// Is this a bug in localAI? Why is the tool call not being parsed correctly?
func getToolCall(input string) string {
	if input[0:2] == "{{" {
		toolCallLength := len(input)
		return input[1 : toolCallLength-1]
	}
	return input
}
