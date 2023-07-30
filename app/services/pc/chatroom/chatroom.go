package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/herman-hang/herman/app/models"
	"github.com/herman-hang/herman/app/repositories"
	"github.com/herman-hang/herman/kernel/core"
)

func List(data map[string]interface{}, ctx *gin.Context) map[string]interface{} {
	users, _ := ctx.Get("pc")

	model := users.(models.Users)
	id, err := repositories.Chatroom().GetChatroomWithMessagesByUserId(model.Id, data["keywords"].(string))
	if err != nil {
		return nil
	}
	core.Log.Debug(err, id)
	return data
}
