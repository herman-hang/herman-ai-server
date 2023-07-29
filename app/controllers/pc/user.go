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

// SendCode 发送验证码
// @param *gin.Context ctx 上下文
func SendCode(ctx *gin.Context) {
	context := app.Request{Context: ctx}
	data := context.Params()
	UserService.SendCode(UserValidate.SendCode(data), ctx)
	context.Json(nil)
}

// UserInfo 获取用户信息
// @param *gin.Context ctx 上下文
func UserInfo(ctx *gin.Context) {
	context := app.Request{Context: ctx}
	data := UserService.UserInfo(ctx)
	context.Json(data)
}

// UserModify 用户信息修改
// @param *gin.Context ctx 上下文
func UserModify(ctx *gin.Context) {
	context := app.Request{Context: ctx}
	data := context.Params()
	UserService.UserModify(UserValidate.ModifyInfo.Check(data), ctx)
	context.Json(nil)
}
