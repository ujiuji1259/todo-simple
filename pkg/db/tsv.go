package db

import (
	"os"
	"encoding/csv"
	"strconv"

	"todo-simple/pkg/todo"
)

type TodoFile struct {
	FilePath string
}

func NewTodoFile(filePath string) (*TodoFile, error) {
	return &TodoFile{
		FilePath: filePath,
	}, nil
}

func (t *TodoFile) Load() (map[todo.TodoId]*todo.TodoItem, error) {
	file, err := os.Open(t.FilePath)
    if err != nil {
		_, err = os.Create(t.FilePath)
		if err != nil {
			return nil, err
		}
		return map[todo.TodoId]*todo.TodoItem{}, nil
    }
    defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	todoMap := make(map[todo.TodoId]*todo.TodoItem)
	for _, v := range rows {
		todoIdInt, err := strconv.Atoi(v[0])
		if err != nil {
			return nil, err
		}

		todoId := todo.TodoId(todoIdInt)
		todoStatus, err := todo.TodoStatusString(v[3])
		if err != nil {
			return nil, err
		}

		todoMap[todoId] = &todo.TodoItem{
			Id: todoId,
			TaskName: v[1],
			Project: v[2],
			Status: todoStatus,
		}
	}
	return todoMap, nil
}

func (t *TodoFile) Save(todoMap map[todo.TodoId]*todo.TodoItem) error {
	file, err := os.Create(t.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'

	for _, v := range todoMap {
		err := writer.Write([]string{
			strconv.Itoa(int(v.Id)),
			v.TaskName,
			v.Project,
			v.Status.String(),
		})
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
