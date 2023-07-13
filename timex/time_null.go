package timex

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func UpdateNullTime(t time.Time) NullTime {
	return NullTime{NullTime: sql.NullTime{Valid: true, Time: t}}
}

func (v *NullTime) Now() NullTime {
	v.Time = time.Now()
	v.Valid = true
	return *v
}

func (v NullTime) Add(time_long time.Duration) NullTime {
	v.Time = v.Time.Add(time_long)
	v.Valid = true
	return v
}

func NullTimeNow() NullTime {
	var v = NullTime{
		NullTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	return v
}
func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time.Format(time.RFC3339))
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	var s *time.Time
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Time = *s
	} else {
		v.Valid = false
	}
	return nil
}

func TimeStringToNullTime(str_time string) NullTime {
	if str_time == "" || str_time == "null" {
		return NullTime{NullTime: sql.NullTime{Valid: false, Time: time.Time{}}}
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str_time, LOC_ZONE) //时区
	if err != nil {
		return NullTime{NullTime: sql.NullTime{Valid: false}}
	}
	return NullTime{NullTime: sql.NullTime{Valid: true, Time: t}}
}

func TimeStringToTime(str_time string) time.Time {
	if str_time == "" || str_time == "null" {
		return time.Time{}
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str_time, LOC_ZONE) //时区
	if err != nil {
		return time.Time{}
	}
	return t
}
func TimeToNullTime(t time.Time) NullTime {
	if t.IsZero() {
		return NullTime{NullTime: sql.NullTime{Valid: false, Time: time.Time{}}}
	}
	return NullTime{NullTime: sql.NullTime{Valid: true, Time: t}}
}

func TimeNullToTime(null_time NullTime) time.Time {
	if null_time.Valid {
		return null_time.Time
	}
	return time.Time{}
}
