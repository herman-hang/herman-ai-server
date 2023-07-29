package pc

import (
	"github.com/gin-gonic/gin"
	FileController "github.com/herman-hang/herman/app/controllers/file"
	PcFileController "github.com/herman-hang/herman/app/controllers/pc"
	UserController "github.com/herman-hang/herman/app/controllers/pc"
)

// Router pc端相关路由
// @param *gin.RouterGroup router 路由组对象
// @return void
func Router(router *gin.RouterGroup) {
	// 用户登录
	router.POST("/login", UserController.Login)
	// 发送手机验证码
	router.POST("/send/code", UserController.SendCode)
	// 获取用户信息
	router.GET("/users", UserController.UserInfo)
	// 用户信息修改
	router.PUT("/users", UserController.UserModify)
	// 文件上传
	router.POST("/files/uploads", PcFileController.UploadFile)
	// 文件下载
	router.GET("/files/download/:id", FileController.DownloadFile)
	// 图片预览
	router.GET("/files/preview/:id", FileController.PreviewFile)
}
