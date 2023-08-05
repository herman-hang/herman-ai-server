package openai

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/application"
	OpenAiService "github.com/herman-hang/herman/application/services/pc/openai"
)

// GPT GPT流式响应
// @param *gin.Context ctx 上下文对象
func GPT(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	OpenAiService.GPT(data, ctx)
}
