package cmd

import (
	"fmt"

	"todo-simple/pkg/manager"
	"todo-simple/pkg/todo"
	"github.com/adrg/xdg"
)

type DeleteCmd struct {
	TodoId todo.TodoId `help:"Todo Id"`
}

func (c *DeleteCmd) Run() error {
	todoFilePath, err := xdg.DataFile("todo-simple/todo.tsv")
	if err != nil {
		return err
	}
	todoManager, err := manager.NewTodoManagerFromFile(todoFilePath)
	if err != nil {
		fmt.Printf("Error loading todo manager: %v\n", err)
		return err
	}

	err = todoManager.Delete(c.TodoId)
	if err != nil {
		return fmt.Errorf("Error deleting todo: %v", err)
	}

	err = todoManager.Save()
	if err != nil {
		return fmt.Errorf("Error saving todo manager: %v", err)
	}
	return nil
}
