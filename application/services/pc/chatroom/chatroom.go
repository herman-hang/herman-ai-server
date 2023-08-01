package chatroom

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	ChatroomConstant "github.com/herman-hang/herman/application/constants/pc/chatroom"
	"github.com/herman-hang/herman/application/models"
	"github.com/herman-hang/herman/application/repositories"
	"github.com/herman-hang/herman/kernel/app"
	"github.com/herman-hang/herman/kernel/core"
	"gorm.io/gorm"
	"sort"
	"time"
)

// Add 添加聊天室
// @param map data 前端请求数据
// @param *gin.Context ctx 上下文
// @return void
func Add(data map[string]interface{}, ctx *gin.Context) {
	users, _ := ctx.Get("pc")
	model := users.(models.Users)

	err := core.Db().Transaction(func(tx *gorm.DB) error {
		// 新增聊天室
		chatroom, err := repositories.Chatroom(tx).Insert(map[string]interface{}{
			"ai_type": data["aiType"],
			"name":    data["name"],
		})
		if err != nil {
			return errors.New(ChatroomConstant.AddFail)
		}
		// 新增关联表
		err = repositories.UserChatroom(tx).Create([]map[string]interface{}{
			{
				"user_id":     model.Id,
				"chatroom_id": chatroom["id"],
			}, {
				"user_id":     0, // 0表示机器人
				"chatroom_id": chatroom["id"],
			},
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
	isExist := repositories.UserChatroom().IsExist(map[string]interface{}{
		"chatroom_id": data["id"],
		"user_id":     model.Id,
	})

	if !isExist {
		panic(ChatroomConstant.DataNotExistMessage)
	}

	err := repositories.Chatroom().Update([]uint{data["id"].(uint)}, map[string]interface{}{
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
	app.Log.Debug(data)
	isExist := repositories.UserChatroom().IsExist(map[string]interface{}{
		"chatroom_id": data["id"],
		"user_id":     model.Id,
	})

	if !isExist {
		panic(ChatroomConstant.DataNotExistMessage)
	}

	err := repositories.Chatroom().Delete([]uint{data["id"].(uint)})
	if err != nil {
		panic(ChatroomConstant.DeleteFail)
	}

}

// List 聊天室列表
// @param map data 前端请求数据
// @param *gin.Context ctx 上下文
// @return map[string]interface{} 返回一个数据集合
func List(data map[string]interface{}, ctx *gin.Context) map[string]interface{} {
	users, _ := ctx.Get("pc")
	model := users.(models.Users)
	ids, err := repositories.UserChatroom().GetChatroomIds(model.Id)

	if err != nil {
		panic(ChatroomConstant.GetDataFail)
	}
	chatroom, err := repositories.Chatroom().Get(ids, data)
	if err != nil {
		panic(ChatroomConstant.GetDataFail)
	}
	list := chatroom["list"].([]map[string]interface{})

	if len(list) > 0 {
		for key, value := range list {
			message, _ := repositories.ChatroomMessages().Last(value["id"].(uint))
			createdAt := value["created_at"]
			delete(list[key], "created_at")
			if len(message) > 0 {
				list[key]["createdAt"] = message["created_at"]
				list[key]["newest"] = message["content"]
			} else {
				list[key]["createdAt"] = createdAt
				list[key]["newest"] = "　"
			}

			if message["sender_id"] != model.PhotoId {
				find, _ := repositories.User().Find(map[string]interface{}{
					"id": message["sender_id"],
				}, []string{"id", "photo_id"})
				list[key]["photoId"] = find["photoId"]
			} else {
				find, _ := repositories.User().Find(map[string]interface{}{
					"id": message["receiver_id"],
				}, []string{"id", "photo_id"})
				list[key]["photoId"] = find["photoId"]
			}
		}
		// 排序（倒叙）
		sort.Slice(list, func(i, j int) bool {
			createdAtOne, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", fmt.Sprintf("%s", list[i]["createdAt"]))
			createdAtTwo, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", fmt.Sprintf("%s", list[j]["createdAt"]))
			return createdAtOne.After(createdAtTwo)
		})
		chatroom["list"] = list
	}

	return chatroom
}
