package models

import (
	"gorm.io/gorm"
	"time"
)

// UserChatroom 用户聊天室结构体
type UserChatroom struct {
	Id         uint           `json:"id" gorm:"column:id;primary_key;comment:用户聊天室ID"`
	UserId     uint           `json:"userId" gorm:"column:user_id;comment:用户ID"`
	ChatroomId uint           `json:"chatroomId" gorm:"column:chatroom_id;comment:聊天室ID"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"column:updated_at;comment:更新时间"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 设置表名
func (UserChatroom) TableName() string {
	return "user_chatroom"
}
