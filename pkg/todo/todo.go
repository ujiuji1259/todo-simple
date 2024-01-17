package todo

type TodoStatus int

const (
	Done TodoStatus = iota
	Todo
	Wait
)

type TodoItem struct {
	Id       string     `db:"id"`
	TaskName string     `db:"task"`
	Project  string     `db:"project"`
	Status   TodoStatus `db:"status"`
}

func RotateStatus(status TodoStatus) TodoStatus {
	switch status {
	case Done:
		return Todo
	case Todo:
		return Done
	default:
		return Todo
	}
}

func statusStrings(statuses []TodoStatus) []string {
	var statusStrings []string
	for _, status := range statuses {
		statusStrings = append(statusStrings, status.String())
	}
	return statusStrings
}
