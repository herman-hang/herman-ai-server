package pc

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	PcConstant "github.com/herman-hang/herman/app/constants/pc"
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
		panic(PcConstant.LoginCodeExpire)
	}
	// 判断验证码是否正确
	if result != data["code"] {
		panic(PcConstant.LoginCodeError)
	}

	user, isExist := repositories.User().UserInfoByPhone(fmt.Sprintf("%s", data["phone"]))
	// 判断用户是否存在
	if !isExist {
		user, err := repositories.User().Insert(map[string]interface{}{
			"user":       data["phone"],
			"password":   utils.HashEncode(generateRandomPassword(10)),
			"phone":      data["phone"],
			"nickname":   fmt.Sprintf("用户%s", utils.GenerateVerificationCode()),
			"loginOutIp": c.ClientIP(),
			"loginOutAt": time.Now().Format("2006-01-02 15:04:05"),
			"loginTotal": 1,
		})
		if err != nil {
			panic(PcConstant.LoginFail)
		}
		return utils.GenerateToken(&utils.Claims{Uid: user["id"].(uint), Guard: "pc"})
	} else {
		err := repositories.User().Update([]uint{data["id"].(uint)}, map[string]interface{}{
			"loginOutIp": c.ClientIP(),
			"loginOutAt": time.Now().Format("2006-01-02 15:04:05"),
			"loginTotal": user["loginTotal"].(int) + 1,
		})
		if err != nil {
			panic(PcConstant.LoginFail)
		}
	}

	// 返回token
	return utils.GenerateToken(&utils.Claims{Uid: user["id"].(uint), Guard: "pc"})
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

// UserInfo 获取用户信息
// @param *gin.Context c 上下文
// @param map[string]interface{} 返回用户信息
func UserInfo(ctx *gin.Context) map[string]interface{} {
	admin, _ := ctx.Get("pc")
	info := admin.(*models.Users)
	user, err := repositories.User().Find(map[string]interface{}{
		"id": info.Id,
	}, []string{"id", "phone", "nickname", "photo"})
	if err != nil {
		panic(PcConstant.GetUserInfoFail)
	}
	user["phone"] = hidePhoneNumberMiddle(user["phone"].(string))

	return user
}
