package application

import (
	"fmt"
	"github.com/herman-hang/herman/application/constants"
	"github.com/herman-hang/herman/kernel/utils"
	"net/http"
)

// Response 响应信息结构体
type Response struct {
	HttpCode int         `json:"-"`
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
}

// Option 定义配置选项函数（关键）
type Option func(*Response)

// C 设置JSON结构状态码
// @param int code 状态码
// @return Option 返回配置选项函数
func C(code int) Option {
	return func(this *Response) {
		this.Code = code
	}
}

// M 设置响应信息
// @param string message 自定义响应信息
// @return Option 返回配置选项函数
func M(message string) Option {
	return func(this *Response) {
		this.Message = message
	}
}

// D 设置响应参数
// @param interface{} data 响应数据
// @return Option 返回配置选项函数
func D(data interface{}) Option {
	return func(this *Response) {
		this.Data = data
	}
}

// H 设置HTTP响应状态码
// @param int HttpCode HTTP状态码，比如：200，500等
// @return Option 返回配置选项函数
func H(HttpCode int) Option {
	return func(this *Response) {
		this.HttpCode = HttpCode
	}
}

// Success 方法一：响应函数
// @param *Gin g 上下文结构体
// @param Option opts 接收多个配置选项函数参数，可以是C，M，D，H
func (r *Request) Success(opts ...Option) {
	defaultResponse := &Response{
		HttpCode: http.StatusOK,
		Code:     http.StatusOK,
		Message:  constants.Success,
		Data:     nil,
	}

	// 依次调用opts函数列表中的函数，为结构体成员赋值
	for _, o := range opts {
		o(defaultResponse)
	}
	// 响应http请求
	r.Context.JSON(defaultResponse.HttpCode, defaultResponse)
	return
}

// Json 方法二：响应函数（所有字段转小驼峰写法）
// @param interface{} data 接收响应参数
// @param args 第一个参数为message，第二个参数为code
func (r *Request) Json(data interface{}, args ...interface{}) {
	var jsonString []byte
	// 将数据转为json格式返回
	camelJson, _ := utils.CamelJSON(data)
	switch len(args) {
	case 0:
		jsonString = []byte(fmt.Sprintf(`{"code":%d,"message":"%s","data":%s}`, http.StatusOK, constants.Success, camelJson))
	case 1:
		jsonString = []byte(fmt.Sprintf(`{"code":%d,"message":"%s","data":%s}`, http.StatusOK, args[0], camelJson))
	case 2:
		jsonString = []byte(fmt.Sprintf(`{"code":%d,"message":"%s","data":%s}`, args[1], args[0], camelJson))
	}
	// 响应http请求
	r.Context.Data(http.StatusOK, "application/json", jsonString)
}
