package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/application"
	ChatroomService "github.com/herman-hang/herman/application/services/pc/chatroom"
	ChatroomValidate "github.com/herman-hang/herman/application/validates/pc/chatroom"
)

// AddChatroom 添加聊天室
// @param *gin.Context ctx 上下文
func AddChatroom(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	ChatroomService.Add(ChatroomValidate.Add.Check(data), ctx)
	context.Json(nil)
}

// ModifyChatroom 修改聊天室
// @param *gin.Context ctx 上下文
func ModifyChatroom(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	ChatroomService.Modify(ChatroomValidate.Modify.Check(data), ctx)
	context.Json(nil)
}

// RemoveChatroom 删除聊天室
// @param *gin.Context ctx 上下文
func RemoveChatroom(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	ChatroomService.Delete(ChatroomValidate.Delete.Check(data), ctx)
	context.Json(nil)
}

// List 聊天室列表
// @param *gin.Context ctx 上下文
func List(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	context.Json(ChatroomService.List(ChatroomValidate.List.Check(data), ctx))
}

// SendMessage 聊天室列表
// @param *gin.Context ctx 上下文
func SendMessage(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	context.Json(ChatroomService.Send(ChatroomValidate.Send.Check(data), ctx))
}

// FindMessages 聊天室消息列表
// @param *gin.Context ctx 上下文
func FindMessages(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	context.Json(ChatroomService.Find(ChatroomValidate.Find.Check(data), ctx))
}
