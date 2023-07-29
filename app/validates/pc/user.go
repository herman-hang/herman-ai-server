package user

import "github.com/herman-hang/herman/app/validates"

// ModifyInfo 用户信息修改验证器
var ModifyInfo = validates.Validates{Validate: ModifyValidate{}}

// ModifyValidate 用户修改验证器
type ModifyValidate struct {
	PhotoId  int    `json:"photoId" validate:"omitempty,number" label:"用户头像ID"`
	Nickname string `json:"nickname" validate:"omitempty" label:"昵称"`
}
