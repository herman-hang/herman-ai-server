package repositories

import (
	"github.com/herman-hang/herman/app/models"
	"github.com/herman-hang/herman/kernel/core"
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

	return &ChatroomMessagesRepository{BaseRepository{Model: new(models.ChatroomMessages), Db: core.Db}}
}
