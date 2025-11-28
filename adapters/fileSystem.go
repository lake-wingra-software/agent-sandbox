package adapters

import (
	"log"
	"os"

	"github.com/openai/openai-go/v3"
)

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(data), nil
}

var ReadFileTool = openai.ChatCompletionToolUnionParam{
	OfFunction: &openai.ChatCompletionFunctionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "read_file",
			Description: openai.String("Read a file at a given path"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"path"},
			},
		},
	},
}

func WriteFile(path string, data string) error {
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(data))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

var WriteFileTool = openai.ChatCompletionToolUnionParam{
	OfFunction: &openai.ChatCompletionFunctionToolParam{
		Function: openai.FunctionDefinitionParam{
			Name:        "write_file",
			Description: openai.String("Write data to a file at a given path"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{
						"type": "string",
					},
					"data": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"path", "data"},
			},
		},
	},
}
