package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/adrg/xdg"
	"todo-simple/pkg/todo"
)

type ListCmd struct {
	Project []string          `help:"Project name"`
	Status  []todo.TodoStatus `help:"Status"`
}

func (c *ListCmd) Run() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	todoPath := filepath.Join(xdg.DataHome, "todo-simple")
	db, err := todo.NewTsvDb(ctx, todoPath)
	if err != nil {
		fmt.Printf("Error initializing todo manager: %v\n", err)
		return err
	}
	items, err := db.ListItems(ctx, c.Project, c.Status)
	if err != nil {
		return fmt.Errorf("failed to list: %w", err)
	}
	renderTodoItems(items)
	return nil
}

func renderTodoItems(todoItems []*todo.TodoItem) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Task Name", "Project Name", "Status"})
	for _, v := range todoItems {
		table.Append([]string{
			v.Id,
			v.TaskName,
			v.Project,
			v.Status.String(),
		})
	}
	table.Render()
}
