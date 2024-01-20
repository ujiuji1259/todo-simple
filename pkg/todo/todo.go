package todo

import (
	"fmt"
	"database/sql/driver"
	"time"

	"github.com/rs/xid"
)

type TodoStatus int

const (
	Done TodoStatus = iota
	Todo
	Wait
)

type TodoItem struct {
	Id         string        `db:"id"`
	TaskName   string        `db:"task"`
	Project    string        `db:"project"`
	Status     TodoStatus    `db:"status"`
	Due        NullTime      `db:"due"`
	Estimation NullDuration  `db:"estimation"`
	StartedAt  NullTime      `db:"started_at"`
	EndedAt    NullTime      `db:"ended_at"`
}

func NewTodoItem(taskName string, projectName string, due NullTime, estimation NullDuration) (*TodoItem, error) {
	guid := xid.New()

	todoStatus, err := TodoStatusString("Todo")
	if err != nil {
		return nil, fmt.Errorf("Error initializing todo status: %w\n", err)
	}

	return &TodoItem{
		Id:         guid.String(),
		TaskName:   taskName,
		Project:    projectName,
		Status:     todoStatus,
		Due: 	    due,
		Estimation: estimation,
		StartedAt:  NullTime{Time: time.Time{}, Valid: false},
		EndedAt:  NullTime{Time: time.Time{}, Valid: false},
	}, nil
}

func StatusStrings(statuses []TodoStatus) []string {
	var statusStrings []string
	for _, status := range statuses {
		statusStrings = append(statusStrings, status.String())
	}
	return statusStrings
}

func (s TodoStatus) Value() (driver.Value, error) {
	return driver.Value(s.String()), nil
}

func (s *TodoStatus) Scan(src interface{}) error {
	todoStatus, err := TodoStatusString(src.(string))
	if err != nil {
		return fmt.Errorf("invalid todo status: %w", err)
	}
	*s = todoStatus
	return nil
}
