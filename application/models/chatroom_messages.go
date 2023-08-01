package models

import (
	"gorm.io/gorm"
	"time"
)

// ChatroomMessages 聊天室消息结构体
type ChatroomMessages struct {
	Id         uint           `json:"id" gorm:"column:id;primary_key;comment:聊天室消息ID"`
	SenderId   uint           `json:"senderId" gorm:"column:sender_id;comment:发送者ID"`
	ReceiverId uint           `json:"receiverId" gorm:"column:receiver_id;comment:接收者ID"`
	Content    string         `json:"content" gorm:"column:content;comment:消息内容"`
	ChatroomId uint           `json:"chatroomId" gorm:"column:chatroom_id;comment:聊天室ID"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
	UpdatedAt  time.Time      `json:"updatedAt" gorm:"column:updated_at;comment:更新时间"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at;index;comment:删除时间"`
}

// TableName 设置表名
func (ChatroomMessages) TableName() string {
	return "chatroom_messages"
}
