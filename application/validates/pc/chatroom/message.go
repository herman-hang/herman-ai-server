package chatroom

import "github.com/herman-hang/herman/application/validates"

// Send 重写验证器结构体，切记不使用引用，而是拷贝
var Send = validates.Validates{Validate: SendValidate{}}

type SendValidate struct {
	ChatroomId uint   `json:"chatroomId" validate:"required,number" label:"聊天室ID"`
	Content    string `json:"content" validate:"required" label:"消息"`
}
