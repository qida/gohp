package sms

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/sms"
)

type QiniuSMS struct {
	SignatureID string
	TemplateID  string
}

var manager *sms.Manager

func InitQiniuSMS(accessKey, secretKey string) {
	fmt.Println("=========InitQiniuSMS===========")
	auth := auth.New(accessKey, secretKey)
	manager = sms.NewManager(auth)
	fmt.Println("============End=================")
}
func SendQiniuSMS(qiniuSMS QiniuSMS, mobile string, code string) (err error) {
	err = CheckRegexMobile(mobile)
	if err != nil {
		return
	}
	if RequestRegLimit <= 0 {
		RequestRegLimit = 0
		err = errors.New("您的请求太过频繁")
		return
	}
	RequestRegLimit--
	if code == "" {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		code = fmt.Sprintf("%04d", rnd.Int31n(10000))
	}
	// SendMessage
	args := sms.MessagesRequest{
		SignatureID: qiniuSMS.SignatureID,
		TemplateID:  qiniuSMS.TemplateID,
		Mobiles:     []string{mobile},
		Parameters: map[string]interface{}{
			"code": code,
		},
	}
	if manager != nil {
		ret, err1 := manager.SendMessage(args)
		if err1 != nil {
			err = err1
			return
		}
		if len(ret.JobID) == 0 {
			err = errors.New("SendMessage() error: The job id cannot be empty")
		}
	} else {
		err = errors.New("manager is nil")
	}
	return
}
