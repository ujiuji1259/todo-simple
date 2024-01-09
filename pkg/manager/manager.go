package manager

import (
	"os"
	"strconv"
	"slices"
	"errors"

	"github.com/olekukonko/tablewriter"

	"todo-simple/pkg/todo"
	"todo-simple/pkg/db"
)

type TodoManager struct {
	TodoMap map[todo.TodoId]*todo.TodoItem
	Db db.Db
}

func NewTodoManager(db db.Db) (*TodoManager, error) {
	todoMap, err := db.Load()
	if err != nil {
		return nil, err
	}
	return &TodoManager{
		TodoMap: todoMap,
		Db: db,
	}, nil
}

func (t *TodoManager) Save() error {
	err := t.Db.Save(t.TodoMap)
	if err != nil {
		return err
	}
	return nil
}

func (t *TodoManager) getNextTodoId() todo.TodoId {
	maxTodoId := todo.TodoId(-1)
	for k := range t.TodoMap {
		if k > maxTodoId {
			maxTodoId = k
		}
	}
	return todo.TodoId(maxTodoId + 1)
}

func (t *TodoManager) Add(taskName string, projectName string) error {
	todoId := t.getNextTodoId()
	todoItem := todo.TodoItem{
		Id: todoId,
		TaskName: taskName,
		Project: projectName,
		Status: todo.Todo,
	}
	t.TodoMap[todoId] = &todoItem
	return nil
}

func (t *TodoManager) Delete(todoId todo.TodoId) error {
	_, ok := t.TodoMap[todoId]
	if !ok {
		return errors.New("TodoId not found")
	}
	delete(t.TodoMap, todoId)
	return nil
}

func (t *TodoManager) List(projects []string, statuses []todo.TodoStatus) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Task Name", "Project Name", "Status"})
	for _, v := range t.TodoMap {
		if len(projects) > 0 && !slices.Contains(projects, v.Project) {
			continue
		}
		if len(statuses) > 0 && !slices.Contains(statuses, v.Status) {
			continue
		}
		table.Append([]string{
			strconv.Itoa(int(v.Id)),
			v.TaskName,
			v.Project,
			v.Status.String(),
		})
	}
	table.Render()
}
