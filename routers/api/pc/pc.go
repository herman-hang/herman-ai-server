package pc

import (
	"github.com/gin-gonic/gin"
	AdminFileController "github.com/herman-hang/herman/application/controllers/admin/file"
	ChatroomController "github.com/herman-hang/herman/application/controllers/pc/chatroom"
	PcFileController "github.com/herman-hang/herman/application/controllers/pc/file"
	OpenAiController "github.com/herman-hang/herman/application/controllers/pc/openai"
	UserController "github.com/herman-hang/herman/application/controllers/pc/user"
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
	router.GET("/users", UserController.FindUser)
	// 获取聊天室列表
	router.GET("/chat/rooms", ChatroomController.List)
	// 添加聊天室
	router.POST("/chat/rooms", ChatroomController.AddChatroom)
	// 修改聊天室
	router.PUT("/chat/rooms", ChatroomController.ModifyChatroom)
	// 删除聊天室
	router.DELETE("/chat/rooms", ChatroomController.RemoveChatroom)
	// 发送消息
	router.POST("/chat/send/messages", ChatroomController.SendMessage)
	// 获取聊天室消息列表
	router.GET("/chat/messages", ChatroomController.FindMessages)
	// GPT流式响应
	router.GET("/chat/gpt4/stream", OpenAiController.GPT4)
	// 用户信息修改
	router.PUT("/users", UserController.ModifyUser)
	// 文件上传
	router.POST("/files/uploads", PcFileController.UploadFile)
	// 文件下载
	router.GET("/files/download/:id", AdminFileController.DownloadFile)
	// 图片预览
	router.GET("/files/preview/:id", AdminFileController.PreviewFile)
}
