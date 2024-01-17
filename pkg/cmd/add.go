package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
)

type AddCmd struct{}

func (c *AddCmd) Run() error {
	taskName, projectName, err := loadTodoFromUI()
	if err != nil {
		return fmt.Errorf("Error loading todo from UI: %w\n", err)
	}

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
