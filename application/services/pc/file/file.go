package file

import (
	"github.com/gin-gonic/gin"
	FileConstant "github.com/herman-hang/herman/application/constants/admin/file"
	"github.com/herman-hang/herman/application/models"
	"github.com/herman-hang/herman/application/repositories"
	"github.com/herman-hang/herman/application/services/admin/file"
	"github.com/herman-hang/herman/kernel/app"
	"mime/multipart"
)

// Upload 文件上传
// @param ctx *gin.Context 上下文
// @param files []*multipart.FileHeader 文件对象切片
// @return existFileInfos []map[string]interface{} 文件切片信息
func Upload(ctx *gin.Context, files []*multipart.FileHeader) (existFileInfos []map[string]interface{}) {
	// 获取登录信息
	info, _ := ctx.Get("pc")
	user := info.(models.Users)
	app.Log.Debug(user)
	// 执行文件上传
	fileInfos, existFileInfos := file.Exec(files, user.Id)
	for _, info := range fileInfos {
		// 保存文件信息
		fileInfo, err := repositories.File().Insert(info)
		if err != nil {
			panic(FileConstant.RecordFileFail)
		}
		existFileInfos = append(existFileInfos, map[string]interface{}{
			"id":       fileInfo["id"],
			"fileName": fileInfo["fileName"],
			"fileType": fileInfo["fileType"],
			"fileExt":  fileInfo["fileExt"],
			"fileSize": fileInfo["fileSize"],
		})
	}
	return existFileInfos
}
