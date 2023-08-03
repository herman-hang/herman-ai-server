package openai

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/kernel/app"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
)

func GPT(ctx *gin.Context) {
	// 创建 OpenAI 客户端
	client := openai.NewClient(app.Config.OpenAi.SecretKey)

	// 创建 ChatCompletion 请求
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 500,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "帮我使用PHP写一个简单算法",
			},
		},
		Stream: true,
	}

	// 设置响应头，指定为application/json
	ctx.Header("Content-Type", "application/json")

	// 创建 ChatCompletion Stream
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer stream.Close()

	// 创建一个通道，用于接收并发处理的结果
	resultChan := make(chan string)
	// 启动Goroutine进行并发处理
	go processStream(ctx, stream, resultChan)
	// 创建一个缓冲区，用于存储所有响应内容
	var buffer bytes.Buffer
	// 从通道中读取响应内容并发送给客户端
	for result := range resultChan {
		buffer.WriteString(result)
		ctx.JSON(http.StatusOK, gin.H{"content": result})
	}
	// 将缓冲区的内容存入数据库
	saveToDatabase(buffer.String())
	// 结束响应
	ctx.Status(http.StatusOK)
}

// 处理Stream流式响应的函数
func processStream(ctx *gin.Context, stream *openai.ChatCompletionStream, resultChan chan<- string) {
	defer close(resultChan)

	// 处理 Stream 响应并将结果发送到通道中
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// Stream完成时结束循环
			break
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return
		}

		// 将响应内容发送到通道中
		resultChan <- response.Choices[0].Delta.Content
	}
}

func saveToDatabase(content string) {
	app.Log.Debug(content)
}
