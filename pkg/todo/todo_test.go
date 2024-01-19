package todo_test

import (
	"testing"

	"github.com/ujiuji1259/todo-simple/pkg/todo"
)

func TestStatus(t *testing.T) {
	cases := map[string]struct{
		statuses []todo.TodoStatus
		expected []string
	}{
		"All case": {[]todo.TodoStatus{todo.Done, todo.Todo, todo.Wait}, []string{"Done", "Todo", "Wait"}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actual := todo.StatusStrings(c.statuses)

			if len(actual) != len(c.expected) {
				t.Fatalf("expected %s, got %s", c.expected, actual)
			}
			for i := range actual {
				if actual[i] != c.expected[i] {
					t.Fatalf("expected %s, got %s", c.expected, actual)
				}
			}
		})
	}
}
