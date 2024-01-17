package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"todo-simple/pkg/todo"
)

type UpdateCmd struct {
	Id     string          `help:"Task Id"`
	Status todo.TodoStatus `help:"Task status to be updated"`
}

func (c *UpdateCmd) Run() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	todoPath := filepath.Join(xdg.DataHome, "todo-simple")
	db, err := todo.NewTsvDb(ctx, todoPath)
	if err != nil {
		fmt.Printf("Error initializing todo manager: %v\n", err)
		return err
	}
	err = db.UpdateStatus(ctx, c.Id, c.Status)
	if err != nil {
		return fmt.Errorf("failed to rotate status %s: %w", c.Id, err)
	}
	return nil
}
