package plugin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type LiveUrl struct {
	AppName string
	PullUrl string
	PushUrl string
	Key     string

	PushRtmp string
	PlayRtmp string
	PlayM3u8 string
	PlayFlv  string

	TimeUpdate time.Time
}

func (this *LiveUrl) NewLiveUrl(app_name string, pull_url string, key string) (err error) {
	this.AppName = app_name
	this.PullUrl = pull_url
	this.Key = key
	return
}
func (this *LiveUrl) RefreshUrl(id_live int) (err error) {
	t := time.Now().Add(24 * time.Hour).Unix()
	url := fmt.Sprintf("/%s/%d-%d-0-0-%s", this.AppName, id_live, t, this.Key)
	rtmp := fmt.Sprintf("/%s/%d-%d-0-0-%s", this.AppName, id_live, t, this.Key)
	flv := fmt.Sprintf("/%s/%d.flv-%d-0-0-%s", this.AppName, id_live, t, this.Key)
	m3u8 := fmt.Sprintf("/%s/%d.m3u8-%d-0-0-%s", this.AppName, id_live, t, this.Key)

	this.PushRtmp = fmt.Sprintf("rtmp://%s/%s/%d?auth_key=%d-0-0-%s", this.PushUrl, this.AppName, id_live, t, md5V(url))
	this.PlayRtmp = fmt.Sprintf("rtmp://%s/%s/%d?auth_key=%d-0-0-%s", this.PullUrl, this.AppName, id_live, t, md5V(rtmp))
	this.PlayM3u8 = fmt.Sprintf("http://%s/%s/%d.m3u8?auth_key=%d-0-0-%s", this.PullUrl, this.AppName, id_live, t, md5V(m3u8))
	this.PlayFlv = fmt.Sprintf("http://%s/%s/%d.flv?auth_key=%d-0-0-%s", this.PullUrl, this.AppName, id_live, t, md5V(flv))
	this.TimeUpdate = time.Now()
	return
}
func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
