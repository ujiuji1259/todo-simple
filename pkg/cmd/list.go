package cmd

import (
	"fmt"

	"todo-simple/pkg/todo"
	"todo-simple/pkg/manager"
	"github.com/adrg/xdg"
)

type ListCmd struct {
	Project []string `help:"Project name"`
	Status []todo.TodoStatus `help:"Status"`
}

func (c *ListCmd) Run() error {
	todoFilePath, err := xdg.DataFile("todo-simple/todo.tsv")
	if err != nil {
		return err
	}
	todoManager, err := manager.NewTodoManagerFromFile(todoFilePath)
	if err != nil {
		fmt.Printf("Error loading todo manager: %v\n", err)
		return err
	}
	todoManager.List(c.Project, c.Status)
	return nil
}
