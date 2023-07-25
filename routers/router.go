package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/app"
	CaptchaController "github.com/herman-hang/herman/app/controllers/captcha"
	middleware "github.com/herman-hang/herman/middlewares"
	"github.com/herman-hang/herman/routers/api/admin"
	"github.com/herman-hang/herman/routers/api/mobile"
	"github.com/herman-hang/herman/servers/settings"
)

// InitRouter 初始化路由
// @param *gin.Engine rootEngine 路由引擎
// @return *gin.Engine 路由引擎
func InitRouter(rootEngine *gin.Engine) *gin.Engine {
	// 测试路由
	rootEngine.GET("/", func(context *gin.Context) {
		response := app.Request{Context: context}
		response.Success(app.D(map[string]interface{}{
			"message": "Welcome to Herman!",
		}))
	})
	// 设置路由前缀
	api := rootEngine.Group(settings.Config.AppPrefix)
	// 获取验证码
	api.GET("/captcha", CaptchaController.GetCaptcha)
	// 检查验证码正确性
	api.POST("/captcha/check", CaptchaController.CheckCaptcha)

	// 移动端模块
	userRouter := api.Group("/user", middleware.Jwt("user"))
	{
		mobile.Router(userRouter)
	}

	// 后台模块
	adminRouter := api.Group("/admin", middleware.Jwt("admin"), middleware.CheckPermission(), middleware.AdminLogger())
	{
		admin.Router(adminRouter)
	}

	// PC端模块
	pcRouter := api.Group("/pc", middleware.Cors())
	{
		admin.Router(pcRouter)
	}

	return rootEngine
}
