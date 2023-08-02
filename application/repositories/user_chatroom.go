package repositories

import (
	"github.com/herman-hang/herman/application/models"
	"github.com/herman-hang/herman/kernel/core"
	utils "github.com/herman-hang/herman/kernel/utils"
	"gorm.io/gorm"
)

// UserChatroomRepository 聊天室用户关联表仓储层
type UserChatroomRepository struct {
	BaseRepository
}

// UserChatroom 实例化聊天室用户关联表仓储层
// @param *gorm.DB tx 事务
// @return AdminRepository 返回聊天室用户关联表仓储层
func UserChatroom(tx ...*gorm.DB) *UserChatroomRepository {
	if len(tx) > 0 && tx[0] != nil {
		return &UserChatroomRepository{BaseRepository{Model: new(models.UserChatroom), Db: tx[0]}}
	}

	return &UserChatroomRepository{BaseRepository{Model: new(models.UserChatroom), Db: core.Db()}}
}

// GetChatroomIds 查询聊天室ID集合
// @param uint userId 用户ID
// @return chatroomIds err 返回数据和一个错误
func (base UserChatroomRepository) GetChatroomIds(userId uint) (chatroomIds []uint, err error) {
	err = base.Db.Model(&models.UserChatroom{}).Where("user_id = ?", userId).Pluck("chatroom_id", &chatroomIds).Error
	if err != nil {
		return nil, err
	}
	return chatroomIds, nil
}

// FindByChatroomId 根据聊天室ID和发送者ID查询接收者ID
// @param uint chatroomId 聊天室ID
// @param uint userId 发送者ID
// @return error map[string]interface{} 返回一个错误信息和数据
func (base UserChatroomRepository) FindByChatroomId(chatroomId uint, userId uint) (error, map[string]interface{}) {
	data := make(map[string]interface{})
	info := make(map[string]interface{})
	err := base.Db.Model(&models.UserChatroom{}).Select([]string{"id", "user_id"}).Where("chatroom_id = ?", chatroomId).Not("user_id = ?", userId).Find(&data).Error
	if err != nil {
		return err, nil
	}
	if len(data) > 0 {
		for k, v := range data {
			// 下划线转为小驼峰
			info[utils.UnderscoreToLowerCamelCase(k)] = v
		}
	}
	return nil, info
}
