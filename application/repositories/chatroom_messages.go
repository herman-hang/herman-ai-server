package repositories

import (
	"github.com/herman-hang/herman/application/constants"
	"github.com/herman-hang/herman/application/models"
	"github.com/herman-hang/herman/kernel/core"
	"github.com/herman-hang/herman/kernel/utils"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

// ChatroomMessagesRepository 聊天室消息表仓储层
type ChatroomMessagesRepository struct {
	BaseRepository
}

// ChatroomMessages 实例化聊天室消息表仓储层
// @param *gorm.DB tx 事务
// @return AdminRepository 返回聊天室消息表仓储层
func ChatroomMessages(tx ...*gorm.DB) *ChatroomMessagesRepository {
	if len(tx) > 0 && tx[0] != nil {
		return &ChatroomMessagesRepository{BaseRepository{Model: new(models.ChatroomMessages), Db: tx[0]}}
	}

	return &ChatroomMessagesRepository{BaseRepository{Model: new(models.ChatroomMessages), Db: core.Db()}}
}

// Last 查询最后一条消息
// @param uint id 聊天室ID
// @return chatroomMessages err 返回数据模型和一个错误
func (base ChatroomMessagesRepository) Last(id uint) (info map[string]interface{}, err error) {
	data := make(map[string]interface{})
	info = make(map[string]interface{})
	err = base.Db.Model(&models.ChatroomMessages{}).Select([]string{"id", "sender_id", "receiver_id", "content", "created_at"}).
		Where("chatroom_id = ?", id).
		Last(&data).Error
	if err != nil {
		return data, err
	}
	if len(data) > 0 {
		for k, v := range data {
			// 下划线转为小驼峰
			info[utils.UnderscoreToLowerCamelCase(k)] = v
		}
	}
	return info, nil

}

// FindByChatroomId 查询聊天室消息表
// @param map[string]interface{} data 查询条件
// @return info err 返回数据和一个错误
func (base ChatroomMessagesRepository) FindByChatroomId(data map[string]interface{}) (info map[string]interface{}, err error) {
	var (
		page    PageInfo
		total   int64
		pageNum int64
		list    []map[string]interface{}
	)
	if len(data) > 0 {
		if err := mapstructure.WeakDecode(data, &page); err != nil {
			panic(constants.MapToStruct)
		}
	}
	// 总条数
	base.Db.Model(&models.ChatroomMessages{}).Count(&total)
	// 计算总页数
	if page.PageSize != 0 && total%page.PageSize != 0 {
		pageNum = total / page.PageSize
		pageNum++
	}
	err = base.Db.Model(&models.ChatroomMessages{}).
		Select([]string{"id", "sender_id", "receiver_id", "content", "created_at"}).
		Where("chatroom_id = ?", data["chatroomId"]).
		Limit(int(page.PageSize)).
		Offset(int((page.Page - 1) * page.PageSize)).
		Find(&list).Error
	if len(list) > 0 {
		for key, value := range list {
			attributes := make(map[string]interface{})
			for index, item := range value {
				// 下划线转为小驼峰
				attributes[utils.UnderscoreToLowerCamelCase(index)] = item
			}
			list[key] = attributes
		}
	}
	info = map[string]interface{}{
		"list":     list,          // 数据
		"total":    total,         // 总条数
		"pageNum":  pageNum,       // 总页数
		"pageSize": page.PageSize, // 每页大小
		"page":     page.Page,     // 当前页码
	}

	return info, nil
}
