package qywx

import (
	"fmt"
	"time"

	"github.com/qida/gohp/httpx"
	// "github.com/beego/beego/v2/client/httplib"
)

type MsgRestult struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	Invaliduser  string `json:"invaliduser"`
	Invalidparty string `json:"invalidparty"`
	Invalidtag   string `json:"invalidtag"`
}
type MsgMarkdown struct {
	Touser   string `json:"touser"`
	Toparty  string `json:"toparty"`
	Totag    string `json:"totag"`
	Msgtype  string `json:"msgtype"`
	Agentid  int    `json:"agentid"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
	EnableDuplicateCheck int `json:"enable_duplicate_check"`
}

type MsgAccessToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type QYWX struct {
	AccessToken string
	TimeExpires time.Time

	CorpId     string
	CorpSecret string
}

//	func NewQywx(corpid string, corpsecret string) (qy *QYWX, err error) {
//		req := httplib.Get("https://qyapi.weixin.qq.com/cgi-bin/gettoken")
//		req.Param("corpid", corpid)
//		req.Param("corpsecret", corpsecret)
//		var msgAccessToken MsgAccessToken
//		err = req.ToJSON(&msgAccessToken)
//		if err != nil {
//			fmt.Printf("获取token出错：%s\r\n", err.Error())
//			return
//		}
//		qy = &QYWX{
//			AccessToken: msgAccessToken.AccessToken,
//			TimeExpires: time.Now().Add(7200),
//			CorpId:      corpid,
//			CorpSecret:  corpsecret,
//		}
//		return
//	}
func NewQywx(corpid string, corpsecret string) (qy *QYWX, _err error) {
	var msgAccessToken MsgAccessToken
	req := map[string]string{
		"corpid":     corpid,
		"corpsecret": corpsecret,
	}
	_err = httpx.NewClientHttp().Debug(true).GetParams(
		"https://qyapi.weixin.qq.com/cgi-bin/gettoken", req, &msgAccessToken, nil)
	if _err != nil {
		fmt.Printf("获取token出错：%+v\r\n", _err)
		return
	}
	if msgAccessToken.Errcode != 0 {
		_err = fmt.Errorf("获取token出错: errcode=%d, errmsg=%s", msgAccessToken.Errcode, msgAccessToken.Errmsg)
		return
	}
	qy = &QYWX{
		AccessToken: msgAccessToken.AccessToken,
		TimeExpires: time.Now().Add(7200),
		CorpId:      corpid,
		CorpSecret:  corpsecret,
	}
	return
}

func (q *QYWX) GetAccessToken() string {
	if time.Since(q.TimeExpires).Seconds() > 0 {
		// token 过期后重新获取，并更新到当前对象缓存
		if nq, err := NewQywx(q.CorpId, q.CorpSecret); err == nil {
			q.AccessToken = nq.AccessToken
			q.TimeExpires = nq.TimeExpires
		}
	}
	return q.AccessToken
}
