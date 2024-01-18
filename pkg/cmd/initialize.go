package cmd

import (
	"path/filepath"
	"fmt"

	"github.com/adrg/xdg"
	"github.com/ujiuji1259/todo-simple/pkg/todo"
)

func getDb() (*todo.TsvDb, error) {
	todoPath := filepath.Join(xdg.DataHome, "todo-simple")
	db, err := todo.NewTsvDb(todoPath)
	if err != nil {
		return nil, fmt.Errorf("Error initializing todo manager: %w\n", err)
	}
	return db, nil
}
