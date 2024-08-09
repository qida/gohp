package timex

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const OTimeFormat = "15:04:05"

type OTime struct {
	time.Time
}

func (t *OTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = OTime{Time: time.Time{}}
		return
	}
	now, err := time.ParseInLocation(`"`+OTimeFormat+`"`, string(data), LOC_ZONE)
	*t = OTime{Time: now}
	return
}

func (t OTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(OTimeFormat))
	return []byte(formatted), nil
}

func (t OTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *OTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = OTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to OTime", v)
}
func (t *OTime) IsNull() bool {
	if t == nil {
		return true
	}
	return t.Time.UnixNano() == 0
}

func (t *OTime) Add(td time.Duration) error {
	if t == nil {
		return fmt.Errorf("time is nil")
	}
	*t = OTime{Time: t.Time.Add(td)}
	return nil
}
