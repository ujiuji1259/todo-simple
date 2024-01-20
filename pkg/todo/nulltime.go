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

func (nt NullTime) String() string {
    if !nt.Valid {
        return ""
    }
    return nt.Time.Format(timeFormat)
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
	case nil:
		nt.Valid = false
		return nil
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


type NullDuration struct {
    Duration  time.Duration
    Valid bool 
}

// Scan implements the Scanner interface.
func (nd *NullDuration) Scan(value interface{}) error {
	switch v := value.(type) {
		case string:
			if v == "" {
				nd.Valid = false
				return nil
			}
			d, err := time.ParseDuration(v)
			if err != nil {
				return fmt.Errorf("failed to parse duration: %w", err)
			}
			nd.Duration = d
			nd.Valid = true
		case nil:
			nd.Valid = false
			return nil
		default:
			return fmt.Errorf("unsupported type: %T", v)
	}
    return nil
}

// Value implements the driver Valuer interface.
func (nd NullDuration) Value() (driver.Value, error) {
    if !nd.Valid {
        return nil, nil
    }
    return driver.Value(nd.Duration.String()), nil
}
