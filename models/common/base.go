package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

type Model struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoPKModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	DUPLICATE_REGEX = regexp.MustCompile(`(?i)\bduplicate\b`)
)

func IsDuplicateError(err error) bool {
	return err != nil && DUPLICATE_REGEX.MatchString(err.Error())
}

// MapJsonColumn map json column in db
// MapJsonColumn 数据库map型json字段
type MapJsonColumn struct {
	V map[string]interface{}
}

func (s MapJsonColumn) Value() (driver.Value, error) {
	b, err := json.Marshal(s.V)
	return string(b), err
}

func (s *MapJsonColumn) Scan(v interface{}) error {
	var err error

	switch vt := v.(type) {
	case string:
		err = json.Unmarshal([]byte(vt), &s.V)
	case []byte:
		if vt[0] == 91 {
			//var data []map[string]interface{}
			//err = json.Unmarshal(vt, &data)
			return errors.New("json parse error")
		} else {
			var data map[string]interface{}
			err = json.Unmarshal(vt, &data)
			*s = MapJsonColumn{data}
		}
	default:
		return errors.New("json parse error")
	}
	return err

}