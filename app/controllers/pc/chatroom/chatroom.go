package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/app"
	ChatroomService "github.com/herman-hang/herman/app/services/pc/chatroom"
	ChatroomValidate "github.com/herman-hang/herman/app/validates/pc/chatroom"
)

func List(ctx *gin.Context) {
	context := app.Request{Context: ctx}
	data := context.Params()

	context.Json(ChatroomService.List(ChatroomValidate.List.Check(data), ctx))
}
