package repositories

import (
	"github.com/herman-hang/herman/app/models"
	"github.com/herman-hang/herman/app/utils"
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

// UpdateByUserIdAndChatroomId 根据用户ID和聊天室ID更新聊天室数据
// @param map[string]interface{} condition 查询条件
// @param map[string]interface{} data 待更新数据
// @return error 返回一个错误信息
func (base ChatroomRepository) UpdateByUserIdAndChatroomId(condition map[string]interface{}, data map[string]interface{}) error {
	var attributes = make(map[string]interface{})
	// 驼峰转下划线
	for k, v := range data {
		k := utils.ToSnakeCase(k)
		attributes[k] = v
	}
	if err := base.Db.Where(condition).Updates(attributes).Error; err != nil {
		return err
	}
	return nil
}

// DeleteByUserIdAndChatroomId 根据用户ID和聊天室ID删除聊天室
// @param map[string]interface{} condition 查询条件
// @return error 返回一个错误信息
func (base ChatroomRepository) DeleteByUserIdAndChatroomId(condition map[string]interface{}) error {
	if err := base.Db.Delete(&base.Model, condition).Error; err != nil {
		return err
	}
	return nil
}
