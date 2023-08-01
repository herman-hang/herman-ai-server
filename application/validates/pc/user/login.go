package user

import (
	"fmt"
	"github.com/herman-hang/herman/application/constants"
	CaptchaConstant "github.com/herman-hang/herman/application/constants/common/captcha"
	"github.com/herman-hang/herman/application/validates"
	utils2 "github.com/herman-hang/herman/kernel/utils"
	"github.com/mitchellh/mapstructure"
)

// CaptchaLoginValidate 用户登录验证结构体
type CaptchaLoginValidate struct {
	Phone               string `json:"phone" validate:"required,len=11" label:"手机号码"`
	Code                string `json:"code" validate:"required,len=6" label:"验证码"`
	CaptchaType         string `json:"captchaType" validate:"required" label:"验证码类型"`
	CaptchaVerification string `json:"captchaVerification" validate:"required" label:"验证码最终校验Token"`
	Token               string `json:"token" validate:"required" label:"验证码Token"`
	PointJson           string `json:"pointJson" validate:"required" label:"验证码PointJson"`
}

// SendCodeValidate 发送验证码验证结构体
type SendCodeValidate struct {
	Phone       string `json:"phone" validate:"required,len=11" label:"手机号码"`
	CaptchaType string `json:"captchaType" validate:"required" label:"验证码类型"`
	Token       string `json:"token" validate:"required" label:"验证码Token"`
	PointJson   string `json:"pointJson" validate:"required" label:"验证码PointJson"`
}

// Login 登录验证器
// @param map[string]interface{} data 待验证数据
// @return toMap 返回验证通过的数据
func Login(data map[string]interface{}) (toMap map[string]interface{}) {
	var login CaptchaLoginValidate
	// map赋值给结构体
	if err := mapstructure.WeakDecode(data, &login); err != nil {
		panic(constants.MapToStruct)
	}

	if err := validates.Validate(login); err != nil {
		panic(err.Error())
	}

	// 验证码验证
	err := utils2.Factory().GetService(fmt.Sprintf("%s", data["captchaType"])).Verification(fmt.Sprintf("%s", data["token"]),
		fmt.Sprintf("%s", data["pointJson"]))
	if err != nil {
		panic(CaptchaConstant.CheckCaptchaError)
	}

	toMap, err = utils2.ToMap(&login, "json")
	if err != nil {
		panic(constants.StructToMap)
	}

	return toMap
}

// SendCode 发送验证码验证器
// @param map[string]interface{} data 待验证数据
// @return map[string]interface{} 返回验证通过的数据
func SendCode(data map[string]interface{}) (toMap map[string]interface{}) {
	var send SendCodeValidate
	// map赋值给结构体
	if err := mapstructure.WeakDecode(data, &send); err != nil {
		panic(constants.MapToStruct)
	}

	if err := validates.Validate(send); err != nil {
		panic(err.Error())
	}

	// 验证码二次验证
	err := utils2.Factory().GetService(fmt.Sprintf("%s", data["captchaType"])).Check(fmt.Sprintf("%s", data["token"]),
		fmt.Sprintf("%s", data["pointJson"]))

	if err != nil {
		panic(CaptchaConstant.CheckCaptchaError)
	}

	toMap, err = utils2.ToMap(&send, "json")
	if err != nil {
		panic(constants.StructToMap)
	}

	return toMap
}
