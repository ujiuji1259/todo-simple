package manager_test

import (
	"testing"
	"reflect"

	"todo-simple/pkg/todo"
	"todo-simple/pkg/db"
	"todo-simple/pkg/manager"
)

func TestAdd(t *testing.T) {
	cases := map[string]struct{
		taskName string
		projectName string
		todos map[todo.TodoId]*todo.TodoItem
		want map[todo.TodoId]*todo.TodoItem
	}{
		"add task to null todos": {
			"task1",
			"project1",
			map[todo.TodoId]*todo.TodoItem{},
			map[todo.TodoId]*todo.TodoItem{0: &todo.TodoItem{ Id: todo.TodoId(1), TaskName: "task1", Project: "project1", Status: todo.Todo}},
		},
		"add task to todos": {
			"task1",
			"project1",
			map[todo.TodoId]*todo.TodoItem{0: &todo.TodoItem{ Id: todo.TodoId(1), TaskName: "task1", Project: "project1", Status: todo.Todo}},
			map[todo.TodoId]*todo.TodoItem{0: &todo.TodoItem{ Id: todo.TodoId(1), TaskName: "task1", Project: "project1", Status: todo.Todo}, 1: &todo.TodoItem{ Id: todo.TodoId(2), TaskName: "task1", Project: "project1", Status: todo.Todo}},
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			todoFile, err := db.NewTodoFile( "todo.tsv")
			if err != nil {
				t.Error(err)
			}
			manager := manager.TodoManager{
				TodoMap: tt.todos,
				Db: todoFile,
			}
			manager.Add(tt.taskName, tt.projectName)
			if err != nil {
				t.Errorf("want nil, but got %v", err)
			}
			if reflect.DeepEqual(tt.want, manager.TodoMap) {
				t.Errorf("want = , but got = ")
			}
		})
	}
}