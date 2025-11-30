package llm

import (
	"agent-sandbox/adapters"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"
)

func TriageParams() *openai.ChatCompletionNewParams {
	return &openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(ClassifierPrompt),
		},
		Model: getModel(),
		Seed:  openai.Int(0),
	}
}

func ChatParams() *openai.ChatCompletionNewParams {

	return &openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(ChatPrompt),
		},
		Model: getModel(),
		Seed:  openai.Int(0),
	}
}

func ToolUserParams() *openai.ChatCompletionNewParams {
	return &openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(ToolUserPrompt),
		},
		Model: getModel(),
		Seed:  openai.Int(0),
		Tools: []openai.ChatCompletionToolUnionParam{
			adapters.ReadFileTool,
			adapters.WriteFileTool,
		},
	}
}

func AppendMessage(params *openai.ChatCompletionNewParams, message openai.ChatCompletionMessageParamUnion) {
	params.Messages = append(params.Messages, message)
}

func getModel() string {
	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if model == "" {
		fmt.Fprintln(os.Stderr, "missing OPENAI_MODEL environment variable")
		os.Exit(1)
	}
	return model
}
