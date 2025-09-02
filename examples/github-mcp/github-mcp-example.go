package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

func main() {
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		log.Fatalln("Please set GITHUB_TOKEN in env")
	}
	model := os.Getenv("INFERENCE_MODEL")
	if model == "" {
		model = openai.ChatModelGPT4o
	}
	log.Printf("Using model %s", model)
	client := openai.NewClient()
	ctx := context.TODO()
	params := responses.ResponseNewParams{
		Model:           model,
		Tools:           []responses.ToolUnionParam{
			{
				OfMcp: &responses.ToolMcpParam{
					ServerLabel: "github",
					ServerURL:   "https://api.githubcopilot.com/mcp/x/repos/readonly",
					Headers:     map[string]string{
						"Authorization": "Bearer " + token,
					},
				},
			},
		},
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("Using the tool provided, summarize the five most recent commits on main for the llama-stack repo owned by llamastack?"),
		},
	}

	resp, err := client.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(resp.OutputText())
}
