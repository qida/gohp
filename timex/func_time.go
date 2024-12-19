package timex

import (
	"errors"
	"fmt"
	"math"
	"time"
)

var LOC_ZONE = time.FixedZone("CST", 8*3600) // 东八区

// 获取零日期 舍去秒 2000-01-01T15:04:00+08:00
func UpdateTime(src *time.Time) (_err error) {
	if !src.IsZero() {
		*src = time.Date(2000, time.January, 1, src.Hour(), src.Minute(), 0, 0, LOC_ZONE)
		return
	}
	_err = errors.New("输入时间不能空")
	return
}

// 获取零时间 2023-01-01T00:00:00+08:00
func UpdateDate(src *time.Time) (_err error) {
	if !src.IsZero() {
		*src = time.Date(src.Year(), src.Month(), src.Day(), 0, 0, 0, 0, LOC_ZONE)
		return
	}
	_err = errors.New("输入时间不能空")
	return
}

// 获取零日期整时 2000-01-01T15:00:00+08:00
func UpdateHour(src *time.Time) (_err error) {
	if !src.IsZero() {
		*src = time.Date(2000, time.January, 1, src.Hour(), 0, 0, 0, LOC_ZONE)
		return
	}
	_err = errors.New("输入时间不能空")
	return
}



// 修改为日期整时 2023-01-02T15:00:00+08:00
func UpdateDateTimeHour(src *time.Time) (_err error) {
	if !src.IsZero() {
		*src = time.Date(src.Year(), src.Month(), src.Day(), src.Hour(), 0, 0, 0, LOC_ZONE)
		return
	}
	_err = errors.New("输入时间不能空")
	return
}


// 获取日期分整数 2023-01-02T15:01:00+08:00
func GetDateTimeMinute(src time.Time) time.Time {
	if !src.IsZero() {
		return time.Date(src.Year(), src.Month(), src.Day(), src.Hour(), src.Minute(), 0, 0, LOC_ZONE)
	}
	return time.Time{}
}

// 获取日期时整数 2023-01-02T15:00:00+08:00
func GetDateTimeHour(src time.Time) time.Time {
	if !src.IsZero() {
		return time.Date(src.Year(), src.Month(), src.Day(), src.Hour(), 0, 0, 0, LOC_ZONE)
	}
	return time.Time{}
}


// 获取零日期整数 2023-01-02T00:00:00+08:00

func GetDate(src time.Time) time.Time {
	if !src.IsZero() {
		return time.Date(src.Year(), src.Month(), src.Day(), 0, 0, 0, 0, LOC_ZONE)
	}
	return time.Time{}
}

// 获取指定日期合成后的时间  2023-01-02T15:04:05+08:00

func GetDateTime(date time.Time, src time.Time) (dst time.Time, _err error) {
	if !src.IsZero() {
		dst = time.Date(
			date.Year(), date.Month(), date.Day(),
			src.Hour(), src.Minute(), src.Second(),
			src.Nanosecond()/1e6, LOC_ZONE)
		return
	}
	_err = errors.New("输入时间不能空")
	return
}

// 倒计时
func GetCountDown(t time.Time) string {
	seconds := int(math.Abs(time.Since(t).Seconds()))
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60

	return fmt.Sprintf("%d小时%d分%d秒", hours, minutes, seconds)
}
