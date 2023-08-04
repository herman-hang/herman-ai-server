package openai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	UserConstant "github.com/herman-hang/herman/application/constants/admin/user"
	ChatroomConstant "github.com/herman-hang/herman/application/constants/pc/chatroom"
	OpenAiConstant "github.com/herman-hang/herman/application/constants/pc/openai"
	"github.com/herman-hang/herman/application/repositories"
	"github.com/herman-hang/herman/kernel/app"
	"github.com/herman-hang/herman/kernel/utils"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GPT(data map[string]interface{}, ctx *gin.Context) {
	app.Log.Debug(data)
	var contentBuffer bytes.Buffer
	// 创建一个通道来传递要发送给客户端的数据字符串
	dataChan := make(chan string)
	// 校验token
	claims := checkLogin(ctx, data["token"].(string), dataChan)
	user, _ := repositories.User().Find(map[string]interface{}{
		"id": claims.Uid,
	}, []string{"id", "photo_id"})
	chatroomId, err := strconv.ParseUint(data["chatroomId"].(string), 10, 64)
	userId := user["id"].(uint)

	err, chatroom := repositories.UserChatroom().FindByChatroomId(uint(chatroomId), userId)
	if err != nil || len(chatroom) == 0 {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": ChatroomConstant.SendMessageFail,
			"data":    nil,
		})
		dataChan <- jsonData
		return
	}
	data["receiverId"] = chatroom["userId"]
	// 当前用户为发送者
	senderInfo, err := addMessage(userId, data["receiverId"].(uint), data["content"].(string), uint(chatroomId), dataChan)
	if err != nil {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": ChatroomConstant.SendMessageFail,
			"data":    nil,
		})
		dataChan <- jsonData
		return
	}

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Transfer-Encoding", "chunked")

	// 启动一个 goroutine，等待数据从通道中传入并发送给客户端
	go func() {
		defer close(dataChan)

		for data := range dataChan {
			_ = sse.Encode(ctx.Writer, sse.Event{
				Data: data,
			})
			// 刷新数据，以通知请求端
			ctx.Writer.Flush()
		}

		defer ctx.Abort()
	}()

	json, _ := utils.MapToJson(map[string]interface{}{
		"code":    http.StatusPartialContent,
		"message": "Me",
		"data": map[string]interface{}{
			"id":         senderInfo["id"],
			"isMe":       true,
			"senderId":   userId,
			"receiverId": chatroom["userId"],
			"photoId":    user["photoId"],
			"content":    data["content"],
			"chatroomId": chatroomId,
			"createdAt":  senderInfo["createdAt"],
		},
	})
	dataChan <- json

	receiverInfo, err := processStream(ctx, contentBuffer, dataChan, data, user, uint(chatroomId))
	if err != nil {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": OpenAiConstant.StreamResponseError,
			"data":    nil,
		})
		dataChan <- jsonData
		return
	}
	// 一次性将内容存入数据库
	err = updateToDatabase(map[string]interface{}{
		"id":            receiverInfo["id"],
		"contentBuffer": contentBuffer.String(),
	}, dataChan)
	if err != nil {
		return
	}
}

// 处理Stream流式响应的函数
func processStream(c *gin.Context,
	contentBuffer bytes.Buffer,
	dataChan chan string,
	data map[string]interface{},
	user map[string]interface{},
	chatroomId uint) (map[string]interface{}, error) {
	client := openai.NewClient(app.Config.OpenAi.SecretKey)
	ctx := context.Background()

	request := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: data["content"].(string),
			},
		},
		Stream: true,
		User:   fmt.Sprintf("%s", data["chatroomId"]),
	}
	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": OpenAiConstant.StreamResponseError,
			"data":    nil,
		})
		dataChan <- jsonData
		return nil, err
	}
	defer stream.Close()

	// 当前用户为接收者
	receiverId := data["receiverId"].(uint)

	receiverInfo, err := addMessage(receiverId, user["id"].(uint), "", chatroomId, dataChan)
	if err != nil {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": ChatroomConstant.SendMessageFail,
			"data":    nil,
		})
		dataChan <- jsonData
		return nil, err
	}
	json, _ := utils.MapToJson(map[string]interface{}{
		"code":    http.StatusPartialContent,
		"message": "GPT",
		"data": map[string]interface{}{
			"id":         receiverInfo["id"],
			"isMe":       false,
			"senderId":   receiverId,
			"receiverId": user["id"].(uint),
			"photoId":    0, // 因为OpenAI没有头像，所以使用0代表使用默认头像
			"content":    data["content"],
			"chatroomId": chatroomId,
			"createdAt":  receiverInfo["createdAt"],
		},
	})
	dataChan <- json

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			jsonData, _ := utils.MapToJson(map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": OpenAiConstant.StreamResponseError,
				"data":    nil,
			})
			dataChan <- jsonData
			return nil, err
		}
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusContinue,
			"message": "slice",
			"data":    response.Choices[0].Delta.Content,
		})

		select {
		case dataChan <- jsonData:
		case <-c.Writer.CloseNotify():
			jsonData, _ := utils.MapToJson(map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": OpenAiConstant.StreamResponseClose,
				"data":    nil,
			})
			dataChan <- jsonData
			return nil, err
		}
		// 将内容存入缓冲区
		contentBuffer.WriteString(response.Choices[0].Delta.Content)
	}

	return receiverInfo, nil
}

func updateToDatabase(data map[string]interface{}, dataChan chan string) error {
	err := repositories.ChatroomMessages().Update([]uint{data["id"].(uint)}, map[string]interface{}{
		"content": data["contentBuffer"],
	})
	if err != nil {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": OpenAiConstant.StreamResponseClose,
			"data":    nil,
		})
		dataChan <- jsonData
		return err
	}
	return nil
}

func addMessage(senderId uint, receiverId uint, content string, chatroomId uint, dataChan chan string) (map[string]interface{}, error) {
	newInfo, err := repositories.ChatroomMessages().Insert(map[string]interface{}{
		"senderId":   senderId,
		"receiverId": receiverId,
		"content":    content,
		"chatroomId": chatroomId,
	})
	if err != nil {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": OpenAiConstant.StreamResponseClose,
			"data":    nil,
		})
		dataChan <- jsonData
		return nil, err
	}
	return newInfo, nil
}

// checkLogin 检查登录
func checkLogin(ctx *gin.Context, token string, dataChan chan string) *utils.Claims {
	parts := strings.SplitN(token, " ", UserConstant.SplitByTwo)
	if !(len(parts) == UserConstant.SplitByTwo && parts[0] == "Bearer") {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": UserConstant.TokenExpires,
			"data":    nil,
		})
		dataChan <- jsonData
	}
	return utils.ParseToken(parts[1], ctx, "pc")
}
