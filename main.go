package main

import (
	"agent-sandbox/adapters"
	"agent-sandbox/llm"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/openai/openai-go/v3"
)

func getCompletionOrPanic(client openai.Client, ctx context.Context, params *openai.ChatCompletionNewParams) *openai.ChatCompletion {
	completion, err := client.Chat.Completions.New(ctx, *params)
	if err != nil {
		panic(err)
	}
	return completion
}

func main() {
	messages := adapters.ParseQuery()
	client := llm.CreateClient()
	ctx := context.Background()

	params := classifyPromptAndCreateParams(messages, client, ctx)

	completion := getCompletionOrPanic(client, ctx, params)

	completion = handleToolCalls(completion, params, client, ctx)

	// Print the model's output to stdout
	log.Println("Got model output:")
	log.Println(completion.Choices[0].Message.Content)
	fmt.Print(completion.Choices[0].Message.Content)
}

func handleToolCalls(completion *openai.ChatCompletion, params *openai.ChatCompletionNewParams, client openai.Client, ctx context.Context) *openai.ChatCompletion {
	toolCalls := completion.Choices[0].Message.ToolCalls

	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "read_file" {
			var args map[string]interface{}
			err := json.Unmarshal([]byte(getToolCall(toolCall.Function.Arguments)), &args)

			if err != nil {
				log.Printf("error unmarshalling arguments: %v", err)
				llm.AppendMessage(params, openai.ToolMessage(err.Error(), toolCall.ID))
				continue
			}
			path := args["path"].(string)

			fileContents, err := adapters.ReadFile(path)
			if err != nil {
				log.Printf("error opening file: %v", err)
				llm.AppendMessage(params, openai.ToolMessage(err.Error(), toolCall.ID))
				continue
			}

			llm.AppendMessage(params, openai.ToolMessage(fileContents, toolCall.ID))
		} else if toolCall.Function.Name == "write_file" {
			var args map[string]interface{}
			err := json.Unmarshal([]byte(getToolCall(toolCall.Function.Arguments)), &args)
			if err != nil {
				log.Printf("error unmarshalling arguments: %v", err)
				llm.AppendMessage(params, openai.ToolMessage(err.Error(), toolCall.ID))
				continue
			}
			log.Print(args)
			path := args["path"].(string)
			data := args["data"].(string)

			err = adapters.WriteFile(path, data)
			if err != nil {
				log.Printf("error writing file: %v", err)
				llm.AppendMessage(params, openai.ToolMessage(err.Error(), toolCall.ID))
			}
		} else {
			message := "Unknown tool: " + toolCall.Function.Name
			log.Printf(message)
			llm.AppendMessage(params, openai.ToolMessage(message, toolCall.ID))
		}

		completion = getCompletionOrPanic(client, ctx, params)
		toolCalls = completion.Choices[0].Message.ToolCalls
	}

	return completion
}

func classifyPromptAndCreateParams(messages []openai.ChatCompletionMessageParamUnion, client openai.Client, ctx context.Context) *openai.ChatCompletionNewParams {
	params := llm.TriageParams()
	llm.AppendMessage(params, messages[len(messages)-1])

	completion := getCompletionOrPanic(client, ctx, params)

	log.Println("Prompt type was: " + completion.Choices[0].Message.Content)
	if completion.Choices[0].Message.Content == "chat" {
		params = llm.ChatParams()
	} else {
		params = llm.ToolUserParams()
	}
	llm.AppendMessages(params, messages)
	return params
}

// Is this a bug in localAI? Why is the tool call not being parsed correctly?
func getToolCall(input string) string {
	if input[0:2] == "{{" {
		toolCallLength := len(input)
		return input[1 : toolCallLength-1]
	}
	return input
}
