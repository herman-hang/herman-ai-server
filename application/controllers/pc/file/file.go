package file

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/application"
	FileConstant "github.com/herman-hang/herman/application/constants/admin/file"
	FileService "github.com/herman-hang/herman/application/services/pc/file"
	FileValidate "github.com/herman-hang/herman/application/validates/admin/file"
)

// UploadFile 上传文件
// @param ctx *gin.Context 上下文
// @return void
func UploadFile(ctx *gin.Context) {
	context := application.Request{Context: ctx}
	files := FileValidate.Check(ctx)
	context.Json(FileService.Upload(ctx, files), FileConstant.UploadSuccess)
}
