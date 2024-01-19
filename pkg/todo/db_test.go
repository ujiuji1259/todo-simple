package todo_test

import (
	"context"
	"time"
	"os"
	"testing"
	"path/filepath"

	"github.com/ujiuji1259/todo-simple/pkg/todo"
)

func TestNewTsvDb(t *testing.T) {
	tmpDir := t.TempDir()

	// only exists directory	 
	todoFile := filepath.Join(tmpDir, "todo-simple")
	err := os.Mkdir(filepath.Join(tmpDir, "todo-simple"), 0755)
	if err != nil {
		t.Fatalf("failed to create todo-simple directory: %s", err)
	}

	// already exists todo file
	todoFileWithTsv := filepath.Join(tmpDir, "todo-simple-with-todo")
	err = os.Mkdir(todoFileWithTsv, 0755)
	if err != nil {
		t.Fatalf("failed to create todo-simple-with-todo directory: %s", err)
	}
	err = os.WriteFile(filepath.Join(todoFileWithTsv, "todo.tsv"), []byte("id\ttask\tproject\tstatus\tdue\testimation\n"), 0644)
	if err != nil {
		t.Fatalf("failed to create todo-simple-with-todo tsv file: %s", err)
	}

	cases := map[string]struct{
		targetDir      string
	}{
		"Create tsv directory": {filepath.Join(tmpDir, "todo-simple-with-no-dir")},
		"Create tsv file": {todoFile},
		"Already exists tsv file": {todoFileWithTsv},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			db, err := todo.NewTsvDb(c.targetDir)
			if err != nil {
				t.Fatalf("failed to create tsv db: %s", err)
			}
			defer db.Close()

			if _, err := os.Stat(c.targetDir); err != nil {
				t.Errorf("failed to create tsv directory: %s", err)
			}

			bytes, err := os.ReadFile(filepath.Join(c.targetDir, "todo.tsv"))
			if err != nil {
				t.Fatalf("failed to read todo.tsv: %s", err)
			}
			expected := "id\ttask\tproject\tstatus\tdue\testimation\n"
			if string(bytes) != expected {
				t.Fatalf("expected %s, got %s", expected, string(bytes))
			}
		})

	}
}

func TestListItems(t *testing.T) {
	tmpDir := t.TempDir()

	cases := map[string]struct{
		body     string
		task     []string
		status	 []todo.TodoStatus
		expected []*todo.TodoItem
		isError  bool
	}{
		"Empty todo.tsv": {
			"id\ttask\tproject\tstatus\tdue\testimation\n", 
			[]string{}, 
			[]todo.TodoStatus{}, 
			[]*todo.TodoItem{},
			false,
		},
		"no filter": {
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tTodo\t\t\n", 
			[]string{}, 
			[]todo.TodoStatus{}, 
			[]*todo.TodoItem{
				&todo.TodoItem{
					Id: "0",
					TaskName: "hoge",
					Project: "hoge",
					Status: todo.Todo,
				},
			},
			false,
		},
		"project filter": {
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tTodo\t\t\nfuga\tfuga\tfuga\tTodo\t\t\n",
			[]string{"hoge"}, 
			[]todo.TodoStatus{}, 
			[]*todo.TodoItem{
				&todo.TodoItem{
					Id: "0",
					TaskName: "hoge",
					Project: "hoge",
					Status: todo.Todo,
				},
			},
			false,
		},
		"status filter": {
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tTodo\t\t\nfuga\tfuga\tfuga\tDone\t\t\n",
			[]string{}, 
			[]todo.TodoStatus{todo.Todo}, 
			[]*todo.TodoItem{
				&todo.TodoItem{
					Id: "0",
					TaskName: "hoge",
					Project: "hoge",
					Status: todo.Todo,
				},
			},
			false,
		},
		"invalid status": {
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tHoge\t\t\n",
			[]string{}, 
			[]todo.TodoStatus{}, 
			[]*todo.TodoItem{},
			true,
		},
		"invalid due": {
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tHoge\t2023-12-12T12\t\n",
			[]string{}, 
			[]todo.TodoStatus{}, 
			[]*todo.TodoItem{},
			true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			todoFile := filepath.Join(tmpDir, "todo.tsv")
			err := os.WriteFile(todoFile, []byte(c.body), 0644)
			if err != nil {
				t.Fatalf("failed to create todo.tsv: %s", err)
			}

			db, err := todo.NewTsvDb(tmpDir)
			if err != nil {
				t.Fatalf("failed to create tsv db: %s", err)
			}
			defer db.Close()


			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			items, err := db.ListItems(ctx, c.task, c.status)
			if c.isError && err == nil {
				t.Fatalf("expected error, got nil")
			} else if !c.isError && err != nil {
				t.Fatalf("failed to list items: %s", err)
			}

			if len(items) != len(c.expected) {
				t.Fatalf("expected %d items, got %d items", len(c.expected), len(items))
			}
			for i, item := range items {
				if item.Id != c.expected[i].Id {
					t.Errorf("expected %s, got %s", c.expected[i].Id, item.Id)
				}
				if item.TaskName != c.expected[i].TaskName {
					t.Errorf("expected %s, got %s", c.expected[i].TaskName, item.TaskName)
				}
				if item.Project != c.expected[i].Project {
					t.Errorf("expected %s, got %s", c.expected[i].Project, item.Project)
				}
				if item.Status != c.expected[i].Status {
					t.Errorf("expected %s, got %s", c.expected[i].Status, item.Status)
				}
				if item.Due != c.expected[i].Due {
					t.Errorf("expected %v, got %v", c.expected[i].Due, item.Due)
				}
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tmpDir := t.TempDir()

	cases := map[string]struct{
		body     string
		task     todo.TodoItem
		expected string
		isError  bool
	}{
		"Add task": {
			"id\ttask\tproject\tstatus\tdue\testimation\n", 
			todo.TodoItem{
				Id: "0",
				TaskName: "hoge",
				Project: "hoge",
				Status: todo.Todo,
			},
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tTodo\t\t\n",
			false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			todoFile := filepath.Join(tmpDir, "todo.tsv")
			err := os.WriteFile(todoFile, []byte(c.body), 0644)
			if err != nil {
				t.Fatalf("failed to create todo.tsv: %s", err)
			}

			db, err := todo.NewTsvDb(tmpDir)
			if err != nil {
				t.Fatalf("failed to create tsv db: %s", err)
			}
			defer db.Close()


			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = db.Add(ctx, c.task)
			if c.isError && err == nil {
				t.Fatalf("expected error, got nil")
			} else if !c.isError && err != nil {
				t.Fatalf("failed to list items: %s", err)
			} else if !c.isError {
				result, err := os.ReadFile(todoFile)
				if err != nil {
					t.Fatalf("failed to read todo.tsv: %s", err)
				}
				if string(result) != c.expected {
					t.Fatalf("expected %s, got %s", c.expected, string(result))
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()

	cases := map[string]struct{
		body     string
		id       string
		expected string
		isError  bool
	}{
		"Delete task": {
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tTodo\t\t\n", 
			"0",
			"id\ttask\tproject\tstatus\tdue\testimation\n",
			false,
		},
		"Not existing task id": {
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tTodo\t\t\n", 
			"1",
			"", 
			true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			todoFile := filepath.Join(tmpDir, "todo.tsv")
			err := os.WriteFile(todoFile, []byte(c.body), 0644)
			if err != nil {
				t.Fatalf("failed to create todo.tsv: %s", err)
			}

			db, err := todo.NewTsvDb(tmpDir)
			if err != nil {
				t.Fatalf("failed to create tsv db: %s", err)
			}
			defer db.Close()


			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = db.Delete(ctx, c.id)
			if c.isError && err == nil {
				t.Fatalf("expected error, got nil")
			} else if !c.isError && err != nil {
				t.Fatalf("failed to list items: %s", err)
			} else if !c.isError {
				result, err := os.ReadFile(todoFile)
				if err != nil {
					t.Fatalf("failed to read todo.tsv: %s", err)
				}
				if string(result) != c.expected {
					t.Fatalf("expected %s, got %s", c.expected, string(result))
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tmpDir := t.TempDir()

	cases := map[string]struct{
		body     string
		id       string
		status   todo.TodoStatus
		expected string
		isError  bool
	}{
		"Change status": {
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tTodo\t\t\n", 
			"0",
			todo.Done,
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tDone\t\t\n", 
			false,
		},
		"same status": {
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tTodo\t\t\n", 
			"0",
			todo.Todo,
			"id\ttask\tproject\tstatus\tdue\testimation\n0\thoge\thoge\tTodo\t\t\n", 
			false,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			todoFile := filepath.Join(tmpDir, "todo.tsv")
			err := os.WriteFile(todoFile, []byte(c.body), 0644)
			if err != nil {
				t.Fatalf("failed to create todo.tsv: %s", err)
			}

			db, err := todo.NewTsvDb(tmpDir)
			if err != nil {
				t.Fatalf("failed to create tsv db: %s", err)
			}
			defer db.Close()


			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = db.UpdateStatus(ctx, c.id, c.status)
			if c.isError && err == nil {
				t.Fatalf("expected error, got nil")
			} else if !c.isError && err != nil {
				t.Fatalf("failed to list items: %s", err)
			} else if !c.isError {
				result, err := os.ReadFile(todoFile)
				if err != nil {
					t.Fatalf("failed to read todo.tsv: %s", err)
				}
				if string(result) != c.expected {
					t.Fatalf("expected %s, got %s", c.expected, string(result))
				}
			}
		})
	}
}