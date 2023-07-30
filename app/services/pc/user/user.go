package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	UserConstant "github.com/herman-hang/herman/app/constants/pc/user"
	"github.com/herman-hang/herman/app/models"
	"github.com/herman-hang/herman/app/repositories"
	"github.com/herman-hang/herman/app/utils"
	"github.com/herman-hang/herman/kernel/core"
	"time"
)

// Login 用户登录
// @param map data 前端请求数据
// @param *gin.Context c 上下文
// @return interface{} 返回一个token值
func Login(data map[string]interface{}, c *gin.Context) interface{} {
	// 设置上下文
	ctx := context.Background()
	// 取出手机验证码
	result, err := core.Redis.Get(ctx, fmt.Sprintf("send_code:%s", c.ClientIP())).Result()
	if err != nil {
		panic(UserConstant.LoginCodeExpire)
	}
	// 判断验证码是否正确
	if result != data["code"] {
		panic(UserConstant.LoginCodeError)
	}
	info, err := repositories.User().Find(map[string]interface{}{
		"phone": data["phone"],
	})
	if err != nil {
		panic(UserConstant.LoginFail)
	}
	// 判断用户是否存在
	if len(info) == 0 {
		info, err := repositories.User().Insert(map[string]interface{}{
			"user":       data["phone"],
			"password":   utils.HashEncode(generateRandomPassword(10)),
			"phone":      data["phone"],
			"nickname":   fmt.Sprintf("用户%s", utils.GenerateVerificationCode()),
			"loginOutIp": c.ClientIP(),
			"loginOutAt": time.Now().Format("2006-01-02 15:04:05"),
			"loginTotal": 1,
		})
		if err != nil {
			panic(UserConstant.LoginFail)
		}
		return utils.GenerateToken(&utils.Claims{Uid: info["id"].(uint), Guard: "pc"})
	} else {
		// 更新用户登录信息
		err := repositories.User().Update([]uint{info["id"].(uint)}, map[string]interface{}{
			"loginOutIp": c.ClientIP(),
			"loginOutAt": time.Now().Format("2006-01-02 15:04:05"),
			"loginTotal": info["loginTotal"].(uint) + 1,
		})
		if err != nil {
			panic(UserConstant.LoginFail)
		}
		// 返回token
		return utils.GenerateToken(&utils.Claims{Uid: info["id"].(uint), Guard: "pc"})
	}
}

// SendCode 发送验证码
// @param map data 前端请求数据
// @return void
func SendCode(data map[string]interface{}, ctx *gin.Context) {
	// 生成验证码
	code := utils.GenerateVerificationCode()
	data["code"] = code
	data["content"] = fmt.Sprintf("【Herman AI】您的验证码是：%s，有效期为5分钟，请不要把验证码泄露给其他人。", code)
	//go exec(data)
	code = "123456"

	// 根据IP地址缓存验证码
	// 设置上下文
	c := context.Background()
	core.Redis.Set(c, fmt.Sprintf("send_code:%s", ctx.ClientIP()), code, time.Minute*5)
}

// Find 获取用户信息
// @param *gin.Context c 上下文
// @param map[string]interface{} 返回用户信息
func Find(ctx *gin.Context) map[string]interface{} {
	users, _ := ctx.Get("pc")

	model := users.(models.Users)
	info, err := repositories.User().Find(map[string]interface{}{
		"id": model.Id,
	}, []string{"id", "phone", "nickname", "photo_id"})
	if err != nil {
		panic(UserConstant.GetUserInfoFail)
	}
	info["phone"] = hidePhoneNumberMiddle(info["phone"].(string))

	return info
}

// Modify 用户信息修改
// @param map[string]interface{} data 前端请求数据
// @param *gin.Context ctx 上下文
// @return map[string]interface{} 返回用户信息
func Modify(data map[string]interface{}, ctx *gin.Context) map[string]interface{} {
	users, _ := ctx.Get("pc")
	info := users.(models.Users)
	err := repositories.User().Update([]uint{info.Id}, map[string]interface{}{
		"nickname": data["nickname"],
		"photoId":  data["photoId"],
	})
	if err != nil {
		panic(UserConstant.ModifyInfoFail)
	}
	return data
}
