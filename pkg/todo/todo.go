package todo

type TodoId int
type TodoStatus int

const (
	Done TodoStatus = iota
	Todo
)

type TodoItem struct {
	Id TodoId
	TaskName string
	Project string
	Status TodoStatus
}
