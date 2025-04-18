package timex

import "time"

const (
	LAYOUT_DATE_TIME = "2006-01-02 15:04:05" // 日期时间格式布局
	LAYOUT_DATE      = "2006-01-02 00:00:00" // 日期格式布局
	LAYOUT_TIME      = "2000-01-01 15:04:05" // 时间格式布局

	DATE_TIME_END = "2099-01-01 00:00:00" // 默认结束时间
)

type TimeEnd time.Time

var TIME_END TimeEnd

func (*TimeEnd) Get() time.Time {
	return _time_end
}

var _time_end, _ = time.ParseInLocation(
	"2006-01-02 15:04:05",
	"2099-01-01 00:00:00",
	time.FixedZone("CST", 8*3600), // 东八区
)
