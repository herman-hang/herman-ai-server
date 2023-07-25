package pc

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	PcConstant "github.com/herman-hang/herman/app/constants/pc"
	"github.com/herman-hang/herman/app/repositories"
	"github.com/herman-hang/herman/app/utils"
	"math/big"
	"time"
)

// Login 用户登录
// @param map data 前端请求数据
// @return interface{} 返回一个token值
func Login(data map[string]interface{}, c *gin.Context) interface{} {
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
