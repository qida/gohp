/*
 * @Author: qida
 * @LastEditors: qida
 */
package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

/*
运满友科技短信
*/

type ManYou struct {
	Account  string
	Password string
}

type VerificationParam struct {
	CodeType  string `json:"codetype"`
	SendMsgID string `json:"sendmsgid"`
}

func SendManYouSMS(man_you ManYou, mobile string) (code string, vparam VerificationParam, err error) {
	err = CheckRegexMobile(mobile)
	if err != nil {
		return
	}
	if RequestRegLimit <= 0 {
		RequestRegLimit = 0
		return "", vparam, errors.New("您的请求太过频繁")
	}
	RequestRegLimit--
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code = fmt.Sprintf("%06d", rnd.Int31n(1000000))
	url := fmt.Sprintf("http://sdk.shmyt.cn:8080/manyousms/sendsms?account=%s&password=%s&mobiles=%s&content=您申请的手机验证码是：%s，请输入后进行验证，谢谢！【携手科技】", man_you.Account, man_you.Password, mobile, code)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(body), &vparam)
	return
}
