package main

import (
	"context"
	"log"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)

func main() {
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
					ServerLabel: "servicenow",
					ServerURL:   "http://127.0.0.1:8000/sse",
					AllowedTools: responses.ToolMcpAllowedToolsUnionParam{
						OfMcpAllowedTools: []string{
							"list_incidents",
						},
					},
				},
			},
		},
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String("Using the tool provided, summarize the five most recent incidents."),
		},
	}

	resp, err := client.Responses.New(ctx, params)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(resp.OutputText())
}
