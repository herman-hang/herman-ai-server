package validates

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/herman-hang/herman/application/constants"
	"github.com/herman-hang/herman/kernel/utils"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// Validates 全局验证器
type Validates struct {
	Validate interface{}
}

// ListValidate 管理员列表验证规则
type ListValidate struct {
	Page     uint   `json:"page" validate:"numeric" label:"页码"`
	PageSize uint   `json:"pageSize" validate:"numeric" label:"每页大小"`
	Keywords string `json:"keywords" validate:"omitempty,max=20" label:"每页大小"`
}

// Check 验证方法
// @param map[string]interface{} data 待验证数据
// @return void
func (base Validates) Check(data map[string]interface{}) (toMap map[string]interface{}) {
	// map赋值给结构体
	if err := mapstructure.WeakDecode(data, &base.Validate); err != nil {
		panic(constants.MapToStruct)
	}
	if err := Validate(base.Validate); err != nil {
		panic(err.Error())
	}

	toMap, err := utils.ToMap(base.Validate, "json")

	if err != nil {
		panic(constants.StructToMap)
	}
	return toMap
}

// Validate 全局model数据验证器
// @param 接收一个待数据验证的结构体
// @return err 返回错误信息
func Validate(dataStruct interface{}) (err error) {
	// 验证
	zhCh := zh.New()
	validate := validator.New()

	// 注册一个函数，获取struct tag里自定义的label作为字段名
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("label")
	})

	uni := ut.New(zhCh)
	trans, _ := uni.GetTranslator("zh")
	// 验证器注册翻译器
	if err = zhTrans.RegisterDefaultTranslations(validate, trans); err != nil {
		return err
	}

	if err = validate.Struct(dataStruct); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(trans))
		}
	}

	return nil
}
