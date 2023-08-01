package models

import (
	"gorm.io/gorm"
	"time"
)

// Chatroom 聊天室结构体
type Chatroom struct {
	Id               uint               `json:"id" gorm:"column:id;primary_key;comment:聊天室ID"`
	AiType           uint8              `json:"aiType" gorm:"column:ai_type;default:1;comment:AI类型（1为GPT,2为绘画）"`
	Name             string             `json:"name" gorm:"column:name;comment:聊天室名称"`
	CreatedAt        time.Time          `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
	UpdatedAt        time.Time          `json:"updatedAt" gorm:"column:updated_at;comment:更新时间"`
	DeletedAt        gorm.DeletedAt     `json:"deletedAt" gorm:"column:deleted_at;index;comment:删除时间"`
	ChatroomMessages []ChatroomMessages `json:"chatroomMessages" gorm:"foreignKey:ChatroomId;references:Id"`
}

// TableName 设置表名
func (Chatroom) TableName() string {
	return "chatroom"
}
