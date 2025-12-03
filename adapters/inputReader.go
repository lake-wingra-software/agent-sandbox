package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/openai/openai-go/v3"
)

const aiAgentUsername = "perry-bot"

type githubComment struct {
	Body string `json:"body"`
	User struct {
		Login string `json:"login"`
	} `json:"user"`
}

func ParseQuery() []openai.ChatCompletionMessageParamUnion {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("error reading stdin: %v", err)
	}
	if len(content) == 0 {
		fmt.Fprintln(os.Stderr, "no input provided on stdin; pipe a prompt, e.g.: echo 'Hello' | go run main.go")
		os.Exit(2)
	}

	var comments []githubComment
	if err := json.Unmarshal(content, &comments); err == nil {
		var messages []openai.ChatCompletionMessageParamUnion
		for _, comment := range comments {
			if comment.User.Login == aiAgentUsername {
				messages = append(messages, openai.AssistantMessage(comment.Body))
			} else {
				messages = append(messages, openai.UserMessage(comment.Body))
			}
		}
		return messages
	} else {
		log.Printf("error parsing JSON: %v. Parsing as string instead.", err)
		return []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(string(content)),
		}
	}
}
