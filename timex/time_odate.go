package timex

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	oDateFormat = "2006-01-02"
)

type ODate struct {
	time.Time
}

func (t *ODate) UnmarshalJSON(data []byte) (_err error) {
	if len(data) == 2 {
		*t = ODate{Time: time.Time{}}
		return
	}
	var tt time.Time
	tt, _err = time.ParseInLocation(`"`+oDateFormat+`"`, string(data), LOC_ZONE)
	if _err != nil {
		tt, _err = time.ParseInLocation(`"`+time.RFC3339+`"`, string(data), LOC_ZONE) // 兼容格式
	}
	*t = ODate{Time: tt}
	return
}

func (t ODate) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(oDateFormat))
	return []byte(formatted), nil
}

func (t ODate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *ODate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = ODate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to ODate", v)
}
