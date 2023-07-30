package repositories

import (
	"github.com/herman-hang/herman/app/models"
	"github.com/herman-hang/herman/kernel/core"
	"gorm.io/gorm"
)

// ChatroomRepository 聊天室表仓储层
type ChatroomRepository struct {
	BaseRepository
}

// Chatroom 实例化聊天室表仓储层
// @param *gorm.DB tx 事务
// @return AdminRepository 返回聊天室表仓储层
func Chatroom(tx ...*gorm.DB) *ChatroomRepository {
	if len(tx) > 0 && tx[0] != nil {
		return &ChatroomRepository{BaseRepository{Model: new(models.Chatroom), Db: tx[0]}}
	}

	return &ChatroomRepository{BaseRepository{Model: new(models.Chatroom), Db: core.Db}}
}

// GetChatroomWithMessagesByUserId 根据用户ID关联查询聊天室和聊天消息
func (base ChatroomRepository) GetChatroomWithMessagesByUserId(userId uint, keywords string) ([]models.Chatroom, error) {
	var chatroom []models.Chatroom
	if err := base.Db.
		Select([]string{
			"chatroom.id", "chatroom.name",
			"chatroom_messages.content", "chatroom_messages.created_at",
		}).
		Where("user_id = ?", userId).
		Where("name like ?", "%"+keywords).
		Preload("ChatroomMessages").
		Find(&chatroom).Error; err != nil {
		return nil, err
	}

	return chatroom, nil
}
