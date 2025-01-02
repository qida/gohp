package timex

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	ODateTimeFormat = "2006-01-02 15:04:05"
)

type ODateTime struct {
	time.Time
}

func (t *ODateTime) UnmarshalJSON(data []byte) (_err error) {
	if len(data) == 2 {
		*t = ODateTime{Time: time.Time{}}
		return
	}
	var tt time.Time
	tt, _err = time.ParseInLocation(`"`+ODateTimeFormat+`"`, string(data), LOC_ZONE)
	if _err != nil {
		tt, _err = time.ParseInLocation(`"`+time.RFC3339+`"`, string(data), LOC_ZONE) // 兼容格式
	}
	*t = ODateTime{Time: tt}
	return
}

func (t ODateTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(ODateTimeFormat))
	return []byte(formatted), nil
}

func (t ODateTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *ODateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = ODateTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
