package stringx

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Range struct {
	Min int
	Max int
}

func (t *Range) UnmarshalJSON(data []byte) (_err error) {
	if len(data) == 2 { //""
		*t = Range{Min: 0, Max: 0}
		return
	}
	parts := strings.Split(strings.Trim(string(data), `"`), "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range format 1 [%v]", parts)
	}

	t.Min, _err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if _err != nil {
		fmt.Printf("错误：%v\r\n", parts)
		return
	}
	t.Max, _err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if _err != nil {
		fmt.Printf("错误：%v\r\n", parts)
		return
	}
	return
}

func (t Range) MarshalJSON() ([]byte, error) {
	// if t.Min == 0 && t.Max == 0 {
	// 	return []byte("0-0"), nil
	// }
	return []byte(fmt.Sprintf("%d-%d", t.Min, t.Max)), nil
}

func (t Range) Value() (driver.Value, error) {
	return fmt.Sprintf("%d-%d", t.Min, t.Max), nil
}

func (t *Range) Scan(v interface{}) (_err error) {
	value, ok := v.(string)
	if ok {
		parts := strings.Split(string(value), "-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid range format 2")
		}
		t.Min, _err = strconv.Atoi(strings.TrimSpace(parts[0]))
		if _err != nil {
			fmt.Printf("错误：%v\r\n", parts)
			return
		}
		t.Max, _err = strconv.Atoi(strings.TrimSpace(parts[1]))
		if _err != nil {
			fmt.Printf("错误：%v\r\n", parts)
			return
		}
		return
	}
	_err = fmt.Errorf("can not convert %v to Range", v)
	return
}
