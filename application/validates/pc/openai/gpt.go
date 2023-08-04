package openai

import "github.com/herman-hang/herman/application/validates"

// GPT 重写验证器结构体，切记不使用引用，而是拷贝
var GPT = validates.Validates{Validate: GPTValidate{}}

// GPTValidate 管理员添加验证规则
type GPTValidate struct {
	ChatroomId uint `json:"chatroomId" validate:"required,number" label:"聊天室ID"`
	Content    uint `json:"content" validate:"required" label:"消息内容"`
}
