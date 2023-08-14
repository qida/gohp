package timex

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const ODateFormat = "2006-01-02"

type ODate struct {
	time.Time
}

func (t *ODate) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = ODate{Time: time.Time{}}
		return
	}
	now, err := time.ParseInLocation(`"`+ODateFormat+`"`, string(data), LOC_ZONE)
	*t = ODate{Time: now}
	return
}

func (t ODate) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(ODateFormat))
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
	return fmt.Errorf("can not convert %v to timestamp", v)
}
