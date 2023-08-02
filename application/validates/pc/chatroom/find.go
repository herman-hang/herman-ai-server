package chatroom

import "github.com/herman-hang/herman/application/validates"

// Find 重写验证器结构体，切记不使用引用，而是拷贝
var Find = validates.Validates{Validate: FindValidate{}}

// FindValidate 添加验证规则
type FindValidate struct {
	ChatroomId uint `json:"chatroomId" validate:"required,number" label:"聊天ID"`
	Page       uint `json:"page" validate:"numeric" label:"页码"`
	PageSize   uint `json:"pageSize" validate:"numeric" label:"每页大小"`
}
