package openai

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/application"
	OpenAiService "github.com/herman-hang/herman/application/services/pc/openai"
)

func GPT(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	data := context.Params()
	OpenAiService.GPT(data, ctx)
}
