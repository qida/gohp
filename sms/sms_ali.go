package sms

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type Ali struct {
	Client  *dysmsapi.Client
	Request *dysmsapi.SendSmsRequest
}

func NewAli(id, sk string) *Ali {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = "UPM"
	request.TemplateCode = "SMS_175475271"
	client, _ := dysmsapi.NewClientWithAccessKey("cn-hangzhou", id, sk)
	return &Ali{Client: client, Request: request}
}
func (a *Ali) Send(mobile string) (code string, err error) {
	err = CheckRegexMobile(mobile)
	if err != nil {
		return
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code = fmt.Sprintf("%04v", rnd.Int31n(10000))
	a.Request.PhoneNumbers = mobile
	a.Request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)
	response, err := a.Client.SendSms(a.Request)
	if err != nil {
		return
	}
	fmt.Printf("response is %#v\n", response)
	return
}
