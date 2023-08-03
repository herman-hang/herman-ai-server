package openai

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/kernel/app"
	"github.com/sashabaranov/go-openai"
)

func GPT4(ctx *gin.Context) {
	client := openai.NewClient(app.Config.OpenAi.SecretKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)
	app.Log.Debug(resp, err)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
