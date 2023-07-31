package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/app"
	ChatroomService "github.com/herman-hang/herman/app/services/pc/chatroom"
	ChatroomValidate "github.com/herman-hang/herman/app/validates/pc/chatroom"
)

// AddChatroom 添加聊天室
// @param *gin.Context ctx 上下文
func AddChatroom(ctx *gin.Context) {
	context := app.Request{Context: ctx}
	data := context.Params()
	ChatroomService.Add(ChatroomValidate.Add.Check(data), ctx)
	context.Json(nil)
}

// ModifyChatroom 修改聊天室
// @param *gin.Context ctx 上下文
func ModifyChatroom(ctx *gin.Context) {
	context := app.Request{Context: ctx}
	data := context.Params()
	ChatroomService.Modify(ChatroomValidate.Modify.Check(data), ctx)
	context.Json(nil)
}

// RemoveChatroom 删除聊天室
// @param *gin.Context ctx 上下文
func RemoveChatroom(ctx *gin.Context) {
	context := app.Request{Context: ctx}
	data := context.Params()
	ChatroomService.Delete(ChatroomValidate.Delete.Check(data), ctx)
	context.Json(nil)
}

func List(ctx *gin.Context) {
	context := app.Request{Context: ctx}
	data := context.Params()

	context.Json(ChatroomService.List(ChatroomValidate.List.Check(data), ctx))
}
