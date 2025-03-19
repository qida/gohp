package sms

import (
	"fmt"
	"math/rand"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

type SMSAli struct {
	client  *dysmsapi.Client
	request *dysmsapi.SendSmsRequest
}

// 初始化短信客户端
func New(accessKeyId, accessKeySecret string) *SMSAli {
	config := &openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
		RegionId:        tea.String("cn-hangzhou"), // 默认新加坡（根据实际情况修改）
	}
	return &SMSAli{client: dysmsapi.NewClient(config)}
}

// 发送短信
func (t *SMSAli) Send(mobile string) (code string, err error) {
	err = CheckRegexMobile(mobile)
	if err != nil {
		return
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code = fmt.Sprintf("%04v", rnd.Int31n(10000))

	request := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(mobile),
		SignName:      tea.String("UPM"),
		TemplateCode:  tea.String("SMS_175475271"),
		TemplateParam: tea.String(fmt.Sprintf(`{"code":"%s"}`, code)),
	}
	response, err := t.client.SendSms(request)
	if err != nil {
		return
	}
	// 检查短信发送结果
	if *response.Body.Code != "OK" {
		err = fmt.Errorf("短信发送失败: %s", *response.Body.Message)
		return
	}
	return
}
