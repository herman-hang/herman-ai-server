package chatroom

import (
	"github.com/herman-hang/herman/app/validates"
)

// Add 重写验证器结构体，切记不使用引用，而是拷贝
var Add = validates.Validates{Validate: AddValidate{}}

// AddValidate 管理员添加验证规则
type AddValidate struct {
	AiType int    `json:"aiType" validate:"required,number,oneof=1 2"`
	Name   string `json:"name" validate:"required,max=20"`
}
