package cmd

import (
	"context"
	"fmt"
	"os"
	"time"
	"encoding/csv"

	"github.com/olekukonko/tablewriter"

	"github.com/ujiuji1259/todo-simple/pkg/todo"
)

type ListCmd struct {
	Project []string          `help:"Project name"`
	Status  []todo.TodoStatus `help:"Status"`
	HumanReadable bool `help:"Human readable"`
}

func (c *ListCmd) Run() error {
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

	items, err := db.ListItems(ctx, c.Project, c.Status)
	if err != nil {
		return fmt.Errorf("failed to list: %w", err)
	}

	if c.HumanReadable {
		renderTodoItemsHumanReadable(items)
	} else {
		err := renderTodoItems(items)
		if err != nil {
			return fmt.Errorf("failed to render: %w", err)
		}
	}
	return nil
}

func renderTodoItems(todoItems []*todo.TodoItem) error {
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = '\t'
	defer writer.Flush()

	err := writer.Write([]string{"Id", "Task Name", "Project Name", "Status", "Due", "Estimation"})
	if err != nil {
		return fmt.Errorf("failed to write header: %w\n", err)
	}
	for _, v := range todoItems {
		var due string
		if v.Due.Valid {
			due = v.Due.Time.String()
		}
		var estimation string
		if v.Estimation.Valid {
			estimation = v.Estimation.Duration.String()
		}

		err := writer.Write([]string{
			v.Id,
			v.TaskName,
			v.Project,
			v.Status.String(),
			due,
			estimation,
		})
		if err != nil {
			return fmt.Errorf("failed to write data: %w\n", err)
		}
	}
	return nil
}

func renderTodoItemsHumanReadable(todoItems []*todo.TodoItem) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Task Name", "Project Name", "Status", "Due", "Estimation"})
	for _, v := range todoItems {
		var due string
		if v.Due.Valid {
			due = v.Due.Time.String()
		}
		var estimation string
		if v.Estimation.Valid {
			estimation = v.Estimation.Duration.String()
		}

		table.Append([]string{
			v.Id,
			v.TaskName,
			v.Project,
			v.Status.String(),
			due,
			estimation,
		})
	}
	table.Render()
}
