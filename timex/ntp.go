package timex

import (
	"errors"
	"time"

	"github.com/beevik/ntp"
)

var NTPs = []string{
	"cn.pool.ntp.org",
	"cn.ntp.org.cn",
	"edu.ntp.org.cn",
	"ntp.aliyun.com",
	"ntp1.aliyun.com",
	"ntp2.aliyun.com",
	"ntp3.aliyun.com",
	"ntp4.aliyun.com",
	"ntp5.aliyun.com",
	"ntp6.aliyun.com",
	"ntp7.aliyun.com",
	"hk.ntp.org.cn",
	"sgp.ntp.org.cn",
	"us.ntp.org.cn",
	"time.sunqida.cn",
}

// data = "2021-10-01"
func TimeNtpCompare(date string) (err error) {
	var time1, timeLast time.Time
	var enNet bool
	for i := 0; i < len(NTPs); i++ {
		time1, err = ntp.Time(NTPs[i])
		if err != nil {
			continue
		} else {
			enNet = true
			// fmt.Printf("Date:%v\r\n", time1)
			break
		}
	}
	timeLast, err = time.Parse("2006-01-02", date)
	if err != nil {
		return
	}
	if !enNet || time1.After(timeLast) {
		err = errors.New("参数有误，请保持网络畅通。")
	}
	return
}
