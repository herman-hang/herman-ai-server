package chatroom

import (
	"errors"
	"github.com/gin-gonic/gin"
	ChatroomConstant "github.com/herman-hang/herman/app/constants/pc/chatroom"
	"github.com/herman-hang/herman/app/models"
	"github.com/herman-hang/herman/app/repositories"
	"github.com/herman-hang/herman/kernel/core"
	"gorm.io/gorm"
)

// Add 添加聊天室
// @param map data 前端请求数据
// @param *gin.Context ctx 上下文
// @return void
func Add(data map[string]interface{}, ctx *gin.Context) {
	users, _ := ctx.Get("pc")
	model := users.(models.Users)

	err := core.Db.Transaction(func(tx *gorm.DB) error {
		// 发送者添加
		_, err := repositories.Chatroom(tx).Insert(map[string]interface{}{
			"ai_type": data["aiType"],
			"name":    data["name"],
			"user_id": model.Id,
		})
		if err != nil {
			return errors.New(ChatroomConstant.AddFail)
		}
		// 接收者添加
		_, err = repositories.Chatroom(tx).Insert(map[string]interface{}{
			"ai_type": data["aiType"],
			"name":    data["name"],
			"user_id": 0,
		})
		if err != nil {
			return errors.New(ChatroomConstant.AddFail)
		}

		return nil
	})
	if err != nil {
		panic(err.Error())
	}
}

// Modify 修改聊天室
// @param map data 前端请求数据
// @param *gin.Context ctx 上下文
// @return void
func Modify(data map[string]interface{}, ctx *gin.Context) {
	users, _ := ctx.Get("pc")
	model := users.(models.Users)
	err := repositories.Chatroom().UpdateByUserIdAndChatroomId(map[string]interface{}{
		"id":      data["id"],
		"user_id": model.Id,
	}, map[string]interface{}{
		"name": data["name"],
	})
	if err != nil {
		panic(ChatroomConstant.ModifyFail)
	}
}

// Delete 删除聊天室
// @param map data 前端请求数据
// @param *gin.Context ctx 上下文
// @return void
func Delete(data map[string]interface{}, ctx *gin.Context) {
	users, _ := ctx.Get("pc")
	model := users.(models.Users)
	err := repositories.Chatroom().DeleteByUserIdAndChatroomId(map[string]interface{}{
		"id":      data["id"],
		"user_id": model.Id,
	})
	if err != nil {
		panic(ChatroomConstant.DeleteFail)
	}
}

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
