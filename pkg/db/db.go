package db

import (
	"todo-simple/pkg/todo"
)

type Db interface {
	Load() (map[todo.TodoId]*todo.TodoItem, error)
	Save(todoMap map[todo.TodoId]*todo.TodoItem) error
}
