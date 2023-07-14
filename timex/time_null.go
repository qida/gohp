package timex

import (
	"database/sql"
	"encoding/json"
	"time"
)

type TimeNull struct {
	sql.NullTime
}

func UpdateTimeNull(t time.Time) TimeNull {
	return TimeNull{NullTime: sql.NullTime{Valid: true, Time: t}}
}

func (v *TimeNull) Now() TimeNull {
	v.Time = time.Now()
	v.Valid = true
	return *v
}

func (v TimeNull) Add(time_long time.Duration) TimeNull {
	v.Time = v.Time.Add(time_long)
	v.Valid = true
	return v
}

func TimeNullNow() TimeNull {
	var v = TimeNull{
		NullTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	return v
}
func (v TimeNull) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time.Format(time.RFC3339))
	} else {
		return json.Marshal(nil)
	}
}

func (v *TimeNull) UnmarshalJSON(data []byte) error {
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

func TimeStringToTimeNull(str_time string, loc_zone *time.Location) TimeNull {
	if str_time == "" || str_time == "null" {
		return TimeNull{NullTime: sql.NullTime{Valid: false, Time: time.Time{}}}
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str_time, loc_zone) //时区
	if err != nil {
		return TimeNull{NullTime: sql.NullTime{Valid: false}}
	}
	return TimeNull{NullTime: sql.NullTime{Valid: true, Time: t}}
}

func TimeStringToTime(str_time string, loc_zone *time.Location) time.Time {
	if str_time == "" || str_time == "null" {
		return time.Time{}
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str_time, loc_zone) //时区
	if err != nil {
		return time.Time{}
	}
	return t
}
func TimeToTimeNull(t time.Time) TimeNull {
	if t.IsZero() {
		return TimeNull{NullTime: sql.NullTime{Valid: false, Time: time.Time{}}}
	}
	return TimeNull{NullTime: sql.NullTime{Valid: true, Time: t}}
}

func TimeNullToTime(null_time TimeNull) time.Time {
	if null_time.Valid {
		return null_time.Time
	}
	return time.Time{}
}
