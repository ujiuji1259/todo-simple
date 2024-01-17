package cmd

import (
	"context"
	"fmt"
	"time"
	"path/filepath"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"todo-simple/pkg/todo"
	"github.com/adrg/xdg"
)

type ListCmd struct {
	Project []string `help:"Project name"`
	Status []todo.TodoStatus `help:"Status"`
}

func (c *ListCmd) Run() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	todoPath := filepath.Join(xdg.DataHome, "todo-simple")
	db, err := todo.NewTsvDb(ctx, todoPath)
	if err != nil {
		fmt.Printf("Error initializing todo manager: %v\n", err)
		return err
	}
	items, err := db.List(ctx)
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
			strconv.Itoa(int(v.Id)),
			v.TaskName,
			v.Project,
			v.Status.String(),
		})
	}
	table.Render()
}
