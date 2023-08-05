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
	"github.com/herman-hang/herman/application/repositories"
	"github.com/herman-hang/herman/kernel/app"
	"github.com/herman-hang/herman/kernel/core"
	"github.com/herman-hang/herman/kernel/utils"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GPT OpenAI对接
// @param map[string]interface{} data 请求参数
// @param *gin.Context ctx 请求上下文
// @return void
func GPT(data map[string]interface{}, ctx *gin.Context) {
	var contentBuffer bytes.Buffer
	// 创建一个通道来传递要发送给客户端的数据字符串
	dataChan := make(chan string)
	err := core.Db().Transaction(func(tx *gorm.DB) error {
		// 校验token
		claims, err := checkLogin(ctx, data["token"].(string), dataChan)
		if err != nil {
			return err
		}
		// 查询用户
		user, err := findUser(claims, data, tx)
		if err != nil {
			return err
		}
		// 获取消息内容
		c := context.Background()
		key := fmt.Sprintf("user:%d-chatroom:%d", claims.Uid, user["chatroomId"].(uint))
		data["content"] = core.Redis().Get(c, key).Val()
		core.Redis().Del(c, key)
		// 查询接收者ID
		err, chatroom := repositories.UserChatroom(tx).FindByChatroomId(user["chatroomId"].(uint), user["id"].(uint))
		if err != nil || len(chatroom) == 0 {
			return errors.New(ChatroomConstant.SendMessageFail)
		}
		data["receiverId"] = chatroom["userId"]

		// 当前用户为发送者
		senderInfo, err := addMessage(user["id"].(uint), data["receiverId"].(uint), data["content"].(string), user["chatroomId"].(uint), tx)
		if err != nil {
			return err
		}

		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")
		ctx.Header("Transfer-Encoding", "chunked")

		// 流式响应
		goroutine(dataChan, ctx)

		err = recordSenderMessage(senderInfo, user, chatroom, data, dataChan)
		if err != nil {
			return err
		}

		receiverInfo, err := processStream(ctx, &contentBuffer, dataChan, data, user, user["chatroomId"].(uint), tx)
		if err != nil {
			return errors.New(ChatroomConstant.SendMessageFail)
		}

		// 一次性将内容存入数据库
		err = updateToDatabase(map[string]interface{}{
			"id":            receiverInfo["id"],
			"contentBuffer": contentBuffer.String(),
		}, tx)
		if err != nil {
			return errors.New(ChatroomConstant.SendMessageFail)
		}

		return nil
	})

	// 统一返回错误并回滚数据库
	if err != nil {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    nil,
		})
		dataChan <- jsonData
		ctx.Abort()
	}
}

// 处理Stream流式响应的函数
func processStream(c *gin.Context,
	contentBuffer *bytes.Buffer,
	dataChan chan string,
	data map[string]interface{},
	user map[string]interface{},
	chatroomId uint,
	tx *gorm.DB,
) (map[string]interface{}, error) {
	client := openai.NewClient(app.Config.OpenAi.SecretKey)
	ctx := context.Background()
	messages, err := repositories.ChatroomMessages(tx).FindNewSexDataByChatroomId(chatroomId)
	if err != nil {
		return nil, errors.New(ChatroomConstant.SendMessageFail)
	}
	var contextMessages []openai.ChatCompletionMessage
	// 添加上下文
	if len(messages) > 0 {
		for _, message := range messages {
			var role string
			if message.SenderId == user["id"].(uint) {
				role = openai.ChatMessageRoleUser
			} else {
				role = openai.ChatMessageRoleAssistant
			}
			newMessage := openai.ChatCompletionMessage{
				Role:    role,
				Content: message.Content,
			}
			contextMessages = append([]openai.ChatCompletionMessage{newMessage}, contextMessages...)
		}
	}
	request := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: NumTokensFromMessages(contextMessages, openai.GPT3Dot5Turbo),
		Messages:  contextMessages,
		Stream:    true,
		User:      data["chatroomId"].(string),
	}
	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return nil, errors.New(ChatroomConstant.SendMessageFail)
	}
	defer stream.Close()

	// 当前用户为接收者
	receiverId := data["receiverId"].(uint)

	receiverInfo, err := addMessage(receiverId, user["id"].(uint), "", chatroomId, tx)
	if err != nil {
		return nil, errors.New(ChatroomConstant.SendMessageFail)
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
			return nil, errors.New(ChatroomConstant.SendMessageFail)
		}
		content := response.Choices[0].Delta.Content
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusContinue,
			"message": "slice",
			"data":    content,
		})

		select {
		case dataChan <- jsonData:
		case <-c.Writer.CloseNotify():
			return nil, err
		}
		// 将内容存入缓冲区
		contentBuffer.WriteString(content)
	}
	return receiverInfo, nil
}

// updateToDatabase 更新数据库
// @param map[string]interface{} data 数据
// @return error 错误信息
func updateToDatabase(data map[string]interface{}, tx *gorm.DB) error {
	err := repositories.ChatroomMessages(tx).Update([]uint{data["id"].(uint)}, map[string]interface{}{
		"content": data["contentBuffer"],
	})
	if err != nil {
		return errors.New(ChatroomConstant.SendMessageFail)
	}
	return nil
}

// addMessage 添加消息
// @param uint senderId 发送者id
// @param uint receiverId 接收者id
// @param string content 内容
// @param uint chatroomId 聊天室id
// @param *gorm DB tx 事务
// @return map[string]interface{} error 新消息，错误信息
func addMessage(senderId uint,
	receiverId uint,
	content string,
	chatroomId uint,
	tx *gorm.DB,
) (map[string]interface{}, error) {
	newInfo, err := repositories.ChatroomMessages(tx).Insert(map[string]interface{}{
		"senderId":   senderId,
		"receiverId": receiverId,
		"content":    content,
		"chatroomId": chatroomId,
	})
	if err != nil {
		return nil, errors.New(ChatroomConstant.SendMessageFail)
	}
	return newInfo, nil
}

// checkLogin 检查登录
// @param *gin.Context ctx 上下文
// @param string token 登录token
// @param chan string dataChan 数据通道
// @return *utils.Claims error 用户信息，错误信息
func checkLogin(ctx *gin.Context, token string, dataChan chan string) (*utils.Claims, error) {
	parts := strings.SplitN(token, " ", UserConstant.SplitByTwo)
	if !(len(parts) == UserConstant.SplitByTwo && parts[0] == "Bearer") {
		jsonData, _ := utils.MapToJson(map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": UserConstant.TokenExpires,
			"data":    nil,
		})
		dataChan <- jsonData
		return nil, errors.New(UserConstant.TokenExpires)
	}
	return utils.ParseToken(parts[1], ctx, "pc"), nil
}

// findUser 查找用户信息
// @param *utils.Claims claims 用户信息
// @param map[string]interface{} data 数据
// @param *gorm DB tx 事务
// @return map[string]interface{} error 数据，错误信息
func findUser(claims *utils.Claims, data map[string]interface{}, tx *gorm.DB) (map[string]interface{}, error) {
	user, err := repositories.User(tx).Find(map[string]interface{}{
		"id": claims.Uid,
	}, []string{"id", "photo_id"})
	if err != nil {
		return nil, err
	}
	chatroomId, err := strconv.ParseUint(data["chatroomId"].(string), 10, 64)
	if err != nil {
		return nil, err
	}
	user["chatroomId"] = uint(chatroomId)
	return user, nil
}

// goroutine 启动一个 goroutine，等待数据从通道中传入并发送给客户端
// @param chan string dataChan 数据通道
// @param *gin.Context ctx 上下文
// @return void
func goroutine(dataChan chan string, ctx *gin.Context) {
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

		ctx.Abort()
	}()
}

// recordSenderMessage 记录发送者消息
// @param map[string]interface{} senderInfo 发送者信息
// @param map[string]interface{} user 用户信息
// @param map[string]interface{} chatroom 聊天室信息
// @param map[string]interface{} data 请求数据
// @param chan string dataChan 数据通道
// @return error 错误信息
func recordSenderMessage(senderInfo map[string]interface{},
	user map[string]interface{},
	chatroom map[string]interface{},
	data map[string]interface{},
	dataChan chan string,
) error {
	json, err := utils.MapToJson(map[string]interface{}{
		"code":    http.StatusPartialContent,
		"message": "Me",
		"data": map[string]interface{}{
			"id":         senderInfo["id"],
			"isMe":       true,
			"senderId":   user["id"].(uint),
			"receiverId": chatroom["userId"],
			"photoId":    user["photoId"],
			"content":    data["content"],
			"chatroomId": user["chatroomId"].(uint),
			"createdAt":  senderInfo["createdAt"],
		},
	})
	dataChan <- json
	if err != nil {
		return err
	}
	return nil
}

// NumTokensFromMessages 聊天消息token计算
// @param []openai ChatCompletionMessage messages 消息
// @param string model 模型
// @return int numTokens token数量
func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return
	}

	var tokensPerMessage, tokensPerName int
	switch model {
	case "gpt-3.5-turbo-0613",
		"gpt-3.5-turbo-16k-0613",
		"gpt-4-0314",
		"gpt-4-32k-0314",
		"gpt-4-0613",
		"gpt-4-32k-0613":
		tokensPerMessage = 3
		tokensPerName = 1
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
		tokensPerName = -1   // if there's a name, the role is omitted
	default:
		if strings.Contains(model, "gpt-3.5-turbo") {
			return NumTokensFromMessages(messages, "gpt-3.5-turbo-0613")
		} else if strings.Contains(model, "gpt-4") {
			return NumTokensFromMessages(messages, "gpt-4-0613")
		} else {
			return
		}
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>
	return numTokens
}
