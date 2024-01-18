package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ujiuji1259/todo-simple/pkg/todo"
)

type UpdateCmd struct {
	Id     string          `help:"Task Id"`
	Status todo.TodoStatus `help:"Task status to be updated"`
}

func (c *UpdateCmd) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := getDb()
	if err != nil {
		fmt.Printf("Error initializing todo manager: %v\n", err)
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	err = db.UpdateStatus(ctx, c.Id, c.Status)
	if err != nil {
		return fmt.Errorf("failed to rotate status %s: %w", c.Id, err)
	}
	return nil
}
