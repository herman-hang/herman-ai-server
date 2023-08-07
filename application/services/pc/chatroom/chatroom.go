package chatroom

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	ChatroomConstant "github.com/herman-hang/herman/application/constants/pc/chatroom"
	"github.com/herman-hang/herman/application/models"
	"github.com/herman-hang/herman/application/repositories"
	"github.com/herman-hang/herman/kernel/core"
	"gorm.io/gorm"
	"sort"
	"time"
)

// Add 添加聊天室
// @param map data 前端请求数据
// @param *gin.Context ctx 上下文
// @return void
func Add(data map[string]interface{}, ctx *gin.Context) map[string]interface{} {
	users, _ := ctx.Get("pc")
	model := users.(models.Users)
	var chatroomId uint
	err := core.Db().Transaction(func(tx *gorm.DB) error {
		// 新增聊天室
		chatroom, err := repositories.Chatroom(tx).Insert(map[string]interface{}{
			"ai_type": data["aiType"],
			"name":    data["name"],
		})
		if err != nil {
			return errors.New(ChatroomConstant.AddFail)
		}
		chatroomId = chatroom["id"].(uint)
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

	return map[string]interface{}{
		"id": chatroomId,
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
			createdAt := value["createdAt"]
			delete(list[key], "createdAt")
			if len(message) > 0 {
				var userId uint
				list[key]["createdAt"] = message["createdAt"]
				list[key]["newest"] = message["content"]
				if message["senderId"] != model.PhotoId {
					userId = message["senderId"].(uint)
				} else {
					userId = message["receiverId"].(uint)
				}
				user, err := repositories.User().Find(map[string]interface{}{
					"id": userId,
				}, []string{"id", "photo_id"})
				if err != nil {
					panic(ChatroomConstant.GetDataFail)
				}
				if len(user) > 0 {
					list[key]["photoId"] = user["photoId"]
				} else {
					list[key]["photoId"] = 0
				}
			} else {
				list[key]["createdAt"] = createdAt
				list[key]["newest"] = nil
				list[key]["photoId"] = 0
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

// Send 发送消息
// @param map data 前端请求数据
// @param *gin.Context ctx 上下文
// @return map[string]interface{} 返回一个数据集合
func Send(data map[string]interface{}, ctx *gin.Context) {
	users, _ := ctx.Get("pc")
	model := users.(models.Users)
	chatroomId := data["chatroomId"].(uint)
	// 设置上下文
	c := context.Background()
	core.Redis().Set(c, fmt.Sprintf("user:%d-chatroom:%d", model.Id, chatroomId), data["content"].(string), time.Second*60)
}

// Find 查找聊天室消息
// @param map data 前端请求数据
// @param *gin.Context ctx 上下文
// @return map[string]interface{} 返回一个数据集合
func Find(data map[string]interface{}, ctx *gin.Context) map[string]interface{} {
	users, _ := ctx.Get("pc")
	model := users.(models.Users)
	chatroomMessages, err := repositories.ChatroomMessages().FindByChatroomId(data)
	if err != nil {
		panic(ChatroomConstant.GetDataFail)
	}
	list := chatroomMessages["list"].([]map[string]interface{})
	for key, value := range list {
		var userId uint
		if value["senderId"] == model.Id {
			list[key]["isMe"] = true
			userId = value["senderId"].(uint)
		} else {
			list[key]["isMe"] = false
			userId = value["receiverId"].(uint)
		}
		user, err := repositories.User().Find(map[string]interface{}{
			"id": userId,
		}, []string{"id", "photo_id"})
		if err != nil {
			panic(ChatroomConstant.GetDataFail)
		}
		if len(user) > 0 {
			list[key]["photoId"] = user["photoId"]
		} else {
			list[key]["photoId"] = 0
		}
	}

	chatroomMessages["list"] = list

	return chatroomMessages
}
