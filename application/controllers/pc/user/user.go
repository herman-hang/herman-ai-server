package user

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/application"
	UserConstant "github.com/herman-hang/herman/application/constants/pc/user"
	UserService "github.com/herman-hang/herman/application/services/pc/user"
	UserValidate "github.com/herman-hang/herman/application/validates/pc/user"
)

// Login 用户登录
// @param *gin.Context ctx 上下文
func Login(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	context.Json(UserService.Login(UserValidate.Login(data), ctx), UserConstant.LoginSuccess)
}

// SendCode 发送验证码
// @param *gin.Context ctx 上下文
func SendCode(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	UserService.SendCode(UserValidate.SendCode(data), ctx)
	context.Json(nil)
}

// FindUser 获取用户信息
// @param *gin.Context ctx 上下文
func FindUser(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := UserService.Find(ctx)
	context.Json(data)
}

// ModifyUser 用户信息修改
// @param *gin.Context ctx 上下文
func ModifyUser(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	UserService.Modify(UserValidate.ModifyInfo.Check(data), ctx)
	context.Json(nil)
}
