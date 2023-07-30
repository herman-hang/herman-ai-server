package user

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	SmsConstant "github.com/herman-hang/herman/app/constants/common/sms"
	PcConstant "github.com/herman-hang/herman/app/constants/pc/user"
	"github.com/herman-hang/herman/kernel/core"
	"github.com/herman-hang/herman/servers/settings"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
)

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

// exec 执行发送
// @param map data 前端请求数据
// @return void
func exec(data map[string]interface{}) {
	hash := md5.New()
	// 将字符串转换为字节数组并计算MD5哈希值
	hash.Write([]byte(settings.Config.Sms.Password))
	password := hash.Sum(nil)
	core.Log.Debug(fmt.Sprintf("%x", password))
	// 发起http请求
	response, err := http.Get(fmt.Sprintf("%ssms?u=%s&p=%s&m=%s&c=%s",
		settings.Config.Sms.Api,
		settings.Config.Sms.User,
		fmt.Sprintf("%x", password),
		data["phone"],
		url.QueryEscape(fmt.Sprintf("%s", data["content"])),
	))
	if err != nil {
		panic(PcConstant.SendCodeFail)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(PcConstant.SendCodeFail)
		}
	}(response.Body)
	bodyBytes, _ := ioutil.ReadAll(response.Body)

	// 转为字符串
	code := string(bodyBytes)
	if code != SmsConstant.SendSuccess {
		if SmsConstant.Status[code] != nil {
			panic(SmsConstant.Status[code])
		} else {
			panic(PcConstant.SendCodeFail)
		}
	}
}

// hidePhoneNumberMiddle 手机中间4位设置为密文
// @param string phoneNumber 手机号码
// @return string 返回密文手机号码
func hidePhoneNumberMiddle(phoneNumber string) string {
	if len(phoneNumber) != 11 {
		return phoneNumber
	}

	// 将手机号码的前三位和后四位保留，中间的四位用 * 替代
	return phoneNumber[:3] + "****" + phoneNumber[7:]
}
