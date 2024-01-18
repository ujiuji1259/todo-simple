package todo

import (
	"fmt"
	"time"
	"database/sql/driver"
)

const (
	timeFormat = "2006-01-02 15:04:05 -0700"
)

// because sql.NullTime with saving string is not supported
type NullTime struct {
    Time  time.Time
    Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		if v == "" {
			nt.Valid = false
			return nil
		}
		t, err := time.Parse(timeFormat, v)
		if err != nil {
			return fmt.Errorf("failed to parse time: %w", err)
		}
		nt.Valid = true
		nt.Time = t
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
    return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
    if !nt.Valid {
        return nil, nil
    }
    return driver.Value(nt.Time.Format(timeFormat)), nil
}
