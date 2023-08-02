package chatroom

import "github.com/herman-hang/herman/application/validates"

// Modify 重写验证器结构体，切记不使用引用，而是拷贝
var Modify = validates.Validates{Validate: ModifyValidate{}}

// ModifyValidate 管理员添加验证规则
type ModifyValidate struct {
	Id   uint   `json:"id" validate:"required,number" label:"聊天ID"`
	Name string `json:"name" validate:"required,max=20" label:"名称"`
}
