package todo_test

import (
	"testing"
	"time"

	"github.com/ujiuji1259/todo-simple/pkg/todo"
)

func TestScanNullTime(t *testing.T) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}
	cases := map[string]struct{
		src 	 any
		expected todo.NullTime
		isError  bool
	}{
		"string": {"2024-01-19 00:00:00 +0900", todo.NullTime{Time: time.Date(2024, 1, 19, 0, 0, 0, 0, jst) ,Valid: true}, false},
		"nil": {nil, todo.NullTime{Time: time.Time{} ,Valid: false}, false},
		"other": {1, todo.NullTime{Time: time.Time{} ,Valid: false}, true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var actual todo.NullTime
			err := actual.Scan(c.src)
			if c.isError && err == nil {
				t.Fatalf("expected error, got nil")
			} else if !c.isError && err != nil {
				t.Fatalf("expected nil, got %v", err)
			} else if !c.isError {
				if actual.Valid != c.expected.Valid {
					t.Fatalf("expected %v, got %v", c.expected.Valid, actual.Valid)
				}

				if c.expected.Valid && !actual.Time.Equal(c.expected.Time) {
					t.Fatalf("expected %v, got %v", c.expected.Time, actual.Time)
				}
			}
		})
	}
}

func TestScanNullDuration(t *testing.T) {
	cases := map[string]struct{
		src 	 any
		expected todo.NullDuration
		isError  bool
	}{
		"string": {"1h30m", todo.NullDuration{Duration: time.Hour + 30 * time.Minute, Valid: true}, false},
		"nil": {nil, todo.NullDuration{Duration: time.Duration(0), Valid: false}, false},
		"other": {"hoge", todo.NullDuration{Duration: time.Duration(0), Valid: false}, true},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			var actual todo.NullDuration
			err := actual.Scan(c.src)
			if c.isError && err == nil {
				t.Fatalf("expected error, got nil")
			} else if !c.isError && err != nil {
				t.Fatalf("expected nil, got %v", err)
			} else if !c.isError {
				if actual.Valid != c.expected.Valid {
					t.Fatalf("expected %v, got %v", c.expected.Valid, actual.Valid)
				}

				if c.expected.Valid && actual.Duration != c.expected.Duration {
					t.Fatalf("expected %v, got %v", c.expected.Duration, actual.Duration)
				}
			}
		})
	}
}
