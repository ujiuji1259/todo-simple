package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"todo-simple/pkg/todo"
)

type DeleteCmd struct {
	Id string `help:"Id"`
}

func (c *DeleteCmd) Run() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	todoPath := filepath.Join(xdg.DataHome, "todo-simple")
	db, err := todo.NewTsvDb(ctx, todoPath)
	if err != nil {
		return fmt.Errorf("Error initializing todo manager: %w\n", err)
	}
	err = db.Delete(ctx, c.Id)
	if err != nil {
		return fmt.Errorf("Error deleting %s: %v\n", c.Id, err)
	}

	return nil
}
