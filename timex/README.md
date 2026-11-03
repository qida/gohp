# timex

时间工具包：提供日期计算与格式化、JSON/数据库时间类型、可空时间、定时任务管理、农历公历转换及 NTP 时间同步。

导入：`import "github.com/qida/gohp/timex"`

## 功能列表

### const.go
- `LAYOUT_DATE_TIME` / `LAYOUT_DATE` / `LAYOUT_TIME`：常用时间布局常量
- `TimeEnd`：默认结束时间，`Get()` 返回东八区 2099-01-01 00:00:00

### func_date.go — 日期计算与格式化
- `GetFirstDateOfMonth(d)` / `GetLastDateOfMonth(d)`：所在月份第一天 / 最后一天（0 点）
- `GetFirstDateOfYear(d)`：所在年第一天；`GetZeroTimeOfDay(d)`：当天 0 点；`GetBetweenDates(start, end)`：两日期之间的日期列表
- `GetFirstDateOfWeek(d, start_sunday)`：所在周第一天（可选周日或周一开头）
- `Date(ti, format)` / `DateToString(ts, format)` / `DateT(t, format)`：格式化 Unix 时间戳（int64/字符串）或 time.Time
- `DateParse(dateString, format)`：按 PHP 风格格式（Y-m-d H:i:s 等）解析日期
- `IsToday(t)` / `LgToday(t)` / `LtToday(t)`：判断是否当天 / 大于今天 / 小于今天
- `ParseWithLocation(timeStr)`：按东八区解析 RFC3339 字符串（失败返回当前时间）；`GetMonthStartAndEnd(date)`：当月起止时间

### func_time.go — 时间裁剪与合成
- `UpdateTime(src)`：改为 2000-01-01 并舍去秒；`UpdateDate(src)`：舍去时分秒
- `UpdateHour(src)`：改为 2000-01-01 整时；`UpdateDateTimeHour(src)`：舍去分秒
- `GetDate(src)` / `GetDateTimeHour(src)` / `GetDateTimeMinute(src)`：取零点 / 整时 / 整分（返回值形式）
- `GetDateTime(date, src)`：将日期与时分秒合成一个时间
- `GetCountDown(t)`：距 t 的倒计时字符串（x小时x分x秒）
- `LOC_ZONE`：东八区（CST +8）

### time_odate.go / time_odatetime.go / time_otime.go — JSON/数据库时间类型
- `ODate`：日期类型（2006-01-02），`NewODate(date)` 构造，支持 JSON 序列化与数据库 Scan/Value
- `ODateTime`：日期时间类型（2006-01-02 15:04:05），支持 JSON 与数据库
- `OTime`：时间类型（15:04:05），`NewOTime(date)` 构造，另含 `IsNull()`、`Add(td)`

### time_null.go — 可空时间
- `TimeNull`：sql.NullTime 封装，JSON 序列化时无效值输出 null
- `TimeNullNow()` / `UpdateTimeNull(t)` / `TimeToTimeNull(t)` / `TimeNullToTime(nt)`：构造与互转
- `TimeStringToTimeNull(s, loc)` / `TimeStringToTime(s, loc)`：按 "2006-01-02 15:04:05" 解析字符串

### task.go — 定时任务管理
- `NewTimerTask()`：创建定时任务管理器（`ITimer` 接口 / `LcTimer` 实现，基于 robfig/cron）
- `AddTaskByFunc` / `AddTaskByJob`：按函数或 Job 接口添加任务
- `StartTask` / `StopTask` / `Remove` / `Clear` / `Close` / `FindCron`：任务的启停、移除与查询

### solar.go — 农历公历转换
- `SolarToLuanr(date)`：公历转农历，返回 "YYYY-MM-DD" 字符串及是否闰月
- `SolarToSimpleLuanr(date)`：公历转农历简单格式（如 "2023年05月初六"）
- `SolarToChineseLuanr(date)`：公历转中文农历（干支生肖年 + 腊月/正月 + 初一）
- `LunarToSolar(date, leapMonthFlag)`：农历转公历

### ntp_client.go / ntp_server.go — NTP 时间
- `NTPs`：内置 NTP 服务器列表（cn.pool.ntp.org、阿里云等）
- `TimeNtpCompare(date)`：校验 NTP 网络时间是否晚于指定日期
- `SyncNtp(server)`：查询 NTP 并更新本机系统时间；`UpdateSystemDate(server, dateTime)`：设置系统时间（Windows/Linux/macOS）
- `ServerNtp(port)`：启动 SNTP 服务器（默认端口 123）

## 使用示例

```go
package main

import (
	"fmt"
	"time"

	"github.com/qida/gohp/timex"
)

func main() {
	now := time.Now()
	fmt.Println(timex.GetZeroTimeOfDay(now))             // 当天 0 点
	fmt.Println(timex.DateT(now, "YYYY-MM-DD HH:mm:ss")) // 自定义格式化
	fmt.Println(timex.SolarToSimpleLuanr("2023-06-22"))  // 农历
}
```

## 备注
- 时区统一使用东八区 `LOC_ZONE`；农历转换支持年份范围 1900-2049。
- 依赖：`github.com/robfig/cron/v3`、`github.com/beevik/ntp`、`github.com/gogf/gf`、`github.com/btfak/sntp`；Windows 下修改系统时间需以管理员身份运行。
