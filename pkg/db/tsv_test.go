package db_test

import (
	"os"
	"path/filepath"
	"testing"
	"reflect"

	"todo-simple/pkg/todo"
	"todo-simple/pkg/db"
)

func TestLoad(t *testing.T) {
	tmpdir := t.TempDir()

	cases := map[string]struct{
		body string
		want map[todo.TodoId]*todo.TodoItem
		expectErr bool
	}{
		"normal case": {
			"1\ttask1\tproject1\tTodo",
			map[todo.TodoId]*todo.TodoItem{0: &todo.TodoItem{ Id: todo.TodoId(1), TaskName: "task1", Project: "project1", Status: todo.Todo}},
			false,
		},
		"invalid todo status": {
			"1\ttask1\tproject1\tHoge",
			nil,
			true,
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			fname := filepath.Join(tmpdir, "todo.tsv")
			err := os.WriteFile(fname, []byte(tt.body), 0666)
			if err != nil {
				t.Fatal(err)
			}
			todoFile, err := db.NewTodoFile(fname)
			if err != nil {
				t.Error(err)
			}
			got, err := todoFile.Load()
			if tt.expectErr && err == nil {
				t.Errorf("want error, but got nil")
			} else if !tt.expectErr {
				if err != nil {
					t.Errorf("want nil, but got %v", err)
				}
				if reflect.DeepEqual(tt.want, got) {
					t.Errorf("want = , but got = ")
				}
			}
		})
	}
}

func TestSave(t *testing.T) {
	tmpdir := t.TempDir()

	cases := map[string]struct{
		todos map[todo.TodoId]*todo.TodoItem
		expected string
	}{
		"normal case": {
			map[todo.TodoId]*todo.TodoItem{ 0: &todo.TodoItem{ Id: todo.TodoId(1), TaskName: "task1", Project: "project1", Status: todo.Todo}},
			"1\ttask1\tproject1\tTodo\n",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			todoFile, err := db.NewTodoFile(filepath.Join(tmpdir, "todo.tsv"))
			if err != nil {
				t.Error(err)
			}
			todoFile.Save(tt.todos)

			bytes, err := os.ReadFile(filepath.Join(tmpdir, "todo.tsv"))
			if err != nil {
				t.Error(err)
			}
			if string(bytes) != tt.expected {
				t.Errorf("want = %s, but got = %s", tt.expected, string(bytes))
			}
		})
	}
}