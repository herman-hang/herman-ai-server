package repositories

import (
	"github.com/herman-hang/herman/app/models"
	"github.com/herman-hang/herman/kernel/core"
	"gorm.io/gorm"
)

// UserChatroomRepository 聊天室用户关联表仓储层
type UserChatroomRepository struct {
	BaseRepository
}

// UserChatroom 实例化聊天室用户关联表仓储层
// @param *gorm.DB tx 事务
// @return AdminRepository 返回聊天室用户关联表仓储层
func UserChatroom(tx ...*gorm.DB) *ChatroomRepository {
	if len(tx) > 0 && tx[0] != nil {
		return &ChatroomRepository{BaseRepository{Model: new(models.UserChatroom), Db: tx[0]}}
	}

	return &ChatroomRepository{BaseRepository{Model: new(models.UserChatroom), Db: core.Db}}
}
