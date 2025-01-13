package timex

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTimeOfDay(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取传入的时间所在年第一天
func GetFirstDateOfYear(d time.Time) time.Time {
	return time.Date(d.Year(), 1, 1, 0, 0, 0, 0, LOC_ZONE)
}

// 获取某一天的0点时间
func GetZeroTimeOfDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, LOC_ZONE)
}

func GetFirstDateOfWeek(d time.Time, start_sunday bool) time.Time {
	var timeFirst time.Time
	var n = int(d.Weekday())
	if !start_sunday {
		if d.Weekday() == 0 {
			n = 7
		}
	}
	timeFirst = d.AddDate(0, 0, -1*(n-1))
	timeFirst = GetZeroTimeOfDay(timeFirst)
	return timeFirst
}

func GetBetweenDates(date_start, date_end time.Time) (d []time.Time) {
	if date_end.Before(date_start) {
		// 如果结束时间小于开始时间，异常
		return
	}
	// 输出日期格式固定
	d = append(d, date_start)
	for {
		date_start = date_start.AddDate(0, 0, 1)
		if date_start == date_end {
			break
		}
		d = append(d, date_start)
	}
	return
}

// Format unix time int64 to string
func Date(ti int64, format string) string {
	t := time.Unix(int64(ti), 0)
	return DateT(t, format)
}

// Format unix time string to string
func DateToString(ts string, format string) string {
	i, _ := strconv.ParseInt(ts, 10, 64)
	return Date(i, format)
}

// Format time.Time struct to string
// MM - month - 01
// M - month - 1, single bit
// DD - day - 02
// D - day 2
// YYYY - year - 2006
// YY - year - 06
// HH - 24 hours - 03
// H - 24 hours - 3
// hh - 12 hours - 03
// h - 12 hours - 3
// mm - minute - 04
// m - minute - 4
// ss - second - 05
// s - second = 5
func DateT(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}

// DateFormat pattern rules.
var datePatterns = []string{
	// year
	"Y", "2006", // A full numeric representation of a year, 4 digits   Examples: 1999 or 2003
	"y", "06", //A two digit representation of a year   Examples: 99 or 03

	// month
	"m", "01", // Numeric representation of a month, with leading zeros 01 through 12
	"n", "1", // Numeric representation of a month, without leading zeros   1 through 12
	"M", "Jan", // A short textual representation of a month, three letters Jan through Dec
	"F", "January", // A full textual representation of a month, such as January or March   January through December

	// day
	"d", "02", // Day of the month, 2 digits with leading zeros 01 to 31
	"j", "2", // Day of the month without leading zeros 1 to 31

	// week
	"D", "Mon", // A textual representation of a day, three letters Mon through Sun
	"l", "Monday", // A full textual representation of the day of the week  Sunday through Saturday

	// time
	"g", "3", // 12-hour format of an hour without leading zeros    1 through 12
	"G", "15", // 24-hour format of an hour without leading zeros   0 through 23
	"h", "03", // 12-hour format of an hour with leading zeros  01 through 12
	"H", "15", // 24-hour format of an hour with leading zeros  00 through 23

	"a", "pm", // Lowercase Ante meridiem and Post meridiem am or pm
	"A", "PM", // Uppercase Ante meridiem and Post meridiem AM or PM

	"i", "04", // Minutes with leading zeros    00 to 59
	"s", "05", // Seconds, with leading zeros   00 through 59

	// time zone
	"T", "MST",
	"P", "-07:00",
	"O", "-0700",

	// RFC 2822
	"r", time.RFC1123Z,
}

// Parse Date use PHP time format.
func DateParse(dateString, format string) (time.Time, error) {
	replacer := strings.NewReplacer(datePatterns...)
	format = replacer.Replace(format)
	return time.ParseInLocation(format, dateString, time.Local)
}

// 判断输入时间是否为当天时间
func IsToday(t time.Time) bool {
	return t.Format("2006-01-02") == time.Now().Format("2006-01-02")
}

// 大于今天
func LgToday(t time.Time) bool {
	return t.After(GetZeroTimeOfDay(t).AddDate(0, 0, 1))
}

// 小于今天
func LtToday(t time.Time) bool {
	return t.Before(GetZeroTimeOfDay(t))
}

// 转为本地日期
func ParseWithLocation(timeStr string) (time.Time, error) {
	lt, err := time.ParseInLocation(time.RFC3339, timeStr, time.FixedZone("CST", 8*3600))
	if err != nil {
		return time.Now(), nil
	}
	return lt, nil
}

func GetMonthStartAndEnd(date time.Time) (time.Time, time.Time) {
	// 获取日期所在月的第一天
	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	// 获取日期所在月的最后一天
	monthEnd := monthStart.AddDate(0, 1, 0)
	return monthStart, monthEnd
}
