package pc

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/app"
	UserService "github.com/herman-hang/herman/app/services/pc"
	UserValidate "github.com/herman-hang/herman/app/validates/pc"
)

// Login 用户登录
// @param *gin.Context ctx 上下文
func Login(ctx *gin.Context) {
	context := app.Request{Context: ctx}
	data := context.Params()
	context.Json(UserService.Login(UserValidate.Login(data), ctx))
}
