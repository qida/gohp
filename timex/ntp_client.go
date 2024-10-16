package timex

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/beevik/ntp"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
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

func SyncNtp(server_ntp string) {
	if server_ntp == "" {
		return
	}
	resp, err := ntp.Query(server_ntp)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	err = resp.Validate()
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	tNow := time.Now()
	fmt.Printf("server_ntp %s, stratum %d, offset %+.6f, delay %.5f\n", server_ntp, resp.Stratum, resp.ClockOffset.Seconds(), resp.RTT.Seconds())
	fmt.Println("Current time:", time.Now())
	if !UpdateSystemDate(server_ntp, time.Now().Add(resp.ClockOffset).Format("2006-01-02 15:04:05")) {
		fmt.Println("更新时间 失败")
	}
	fmt.Println(time.Since(tNow).Seconds())
	fmt.Println("Current time:", time.Now())
}

func UpdateSystemDate(server_ntp, dateTime string) bool {
	system := runtime.GOOS
	switch system {
	case "windows":
		{
			_, err1 := gproc.ShellExec(`date  ` + gstr.Split(dateTime, " ")[0])
			_, err2 := gproc.ShellExec(`time  ` + gstr.Split(dateTime, " ")[1])
			if err1 != nil && err2 != nil {
				fmt.Println("更新系统时间错误:请用管理员身份启动程序!")
				return false
			}
			return true
		}
	case "linux":
		{
			_, err1 := gproc.ShellExec(`date -s "` + dateTime + `"`)
			if err1 != nil {
				fmt.Printf("更新系统时间错误1 [%s][%s]:", dateTime, err1.Error())
				return false
			}
			_, err1 = gproc.ShellExec(`ntpdate -u ` + server_ntp)
			if err1 != nil {
				fmt.Println("更新系统时间错误2:", err1.Error())
				return false
			}
			return true
		}
	case "darwin":
		{
			_, err1 := gproc.ShellExec(`date -s "` + dateTime + `"`)
			if err1 != nil {
				fmt.Println("更新系统时间错误:", err1.Error())
				return false
			}
			return true
		}
	}
	return false
}
