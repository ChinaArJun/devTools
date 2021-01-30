package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

const (
	timeFormart = "2006-01-02 15:04:05"
	dateFormat  = "2006-01-02"
	zone        = "Asia/Shanghai"
)

// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(timeFormart))
	return []byte(formatted), nil
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) UnixStr() string {
	unix := t.Unix()
	if unix < 0 {
		return "0"
	}
	return strconv.FormatInt(unix, 10)
}

// UnmarshalJSON implements json unmarshal interface.
func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == `""` {
		return
	}
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	if err != nil {
		now, err = time.ParseInLocation(`"`+dateFormat+`"`, string(data), time.Local)
	}
	t.Time = time.Time(now)
	return
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
