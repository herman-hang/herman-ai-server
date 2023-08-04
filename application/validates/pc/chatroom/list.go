package chatroom

import "github.com/herman-hang/herman/application/validates"

// List 重写验证器结构体，切记不使用引用，而是拷贝
var List = validates.Validates{Validate: ListValidate{}}

// ListValidate 列表验证器
type ListValidate struct {
	AiType   int    `json:"aiType" validate:"required,number,oneof=1 2" label:"AI类型"`
	Page     uint   `json:"page" validate:"numeric" label:"页码"`
	PageSize uint   `json:"pageSize" validate:"numeric" label:"每页大小"`
	Keywords string `json:"keywords" validate:"omitempty,max=20" label:"每页大小"`
}
