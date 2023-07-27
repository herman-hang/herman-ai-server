package pc

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	PcConstant "github.com/herman-hang/herman/app/constants/pc"
	"github.com/herman-hang/herman/app/repositories"
	"github.com/herman-hang/herman/app/utils"
	"github.com/herman-hang/herman/jobs"
	"github.com/herman-hang/herman/kernel/core"
	"math/big"
	"time"
)

// Login 用户登录
// @param map data 前端请求数据
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
			"nickname":   fmt.Sprintf("用户%d", data["id"]),
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

// generateRandomPassword 生成随机密码
// @param int length 密码长度
// @return string 返回一个随机密码
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)
	maxIdx := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			panic(err)
		}
		password[i] = charset[idx.Int64()]
	}

	return string(password)
}

// SendCode 发送验证码
// @param map data 前端请求数据
// @return void
func SendCode(data map[string]interface{}, ctx *gin.Context) {
	// 生成验证码
	code := utils.GenerateVerificationCode()
	// 发送短信
	jobs.Dispatch(map[string]interface{}{
		"topic": "sendSms",
		"data": map[string]interface{}{
			"mobile":  data["phone"],
			"content": fmt.Sprintf("【Herman AI】您的验证码是：%s。请不要把验证码泄露给其他人。", code),
		},
	}, jobs.SendSms)
	// 根据IP地址缓存验证码
	// 设置上下文
	c := context.Background()

	core.Redis.Set(c, fmt.Sprintf("send_code:%s", ctx.ClientIP()), code, time.Minute)
}
