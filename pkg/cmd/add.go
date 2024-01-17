package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/manifoldco/promptui"
	"todo-simple/pkg/todo"
)

type AddCmd struct{}

func (c *AddCmd) Run() error {
	taskName, projectName, err := loadTodoFromUI()
	if err != nil {
		return fmt.Errorf("Error loading todo from UI: %w\n", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	todoPath := filepath.Join(xdg.DataHome, "todo-simple")
	db, err := todo.NewTsvDb(ctx, todoPath)
	if err != nil {
		return fmt.Errorf("Error initializing todo manager: %w\n", err)
	}
	err = db.Add(ctx, taskName, projectName)
	if err != nil {
		return fmt.Errorf("Error loading todo manager: %v\n", err)
	}

	return nil
}

func loadTodoFromUI() (string, string, error) {
	taskNamePrompt := promptui.Prompt{
		Label: "Task Name",
	}
	taskName, err := taskNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", "", err
	}

	projectNamePrompt := promptui.Prompt{
		Label: "Project Name",
	}
	projectName, err := projectNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", "", err
	}

	return taskName, projectName, nil
}
