package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"

	"todo-simple/pkg/todo"
)

type AddCmd struct{}

func (c *AddCmd) Run() error {
	taskName, projectName, due, err := loadTodoFromUI()
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

	newTodoItem, err := todo.NewTodoItem(taskName, projectName, due)
	if err != nil {
		return fmt.Errorf("Error initializing todo item: %w\n", err)
	}
	err = db.Add(ctx, *newTodoItem)
	if err != nil {
		return fmt.Errorf("Error loading todo manager: %v\n", err)
	}

	return nil
}

func loadTodoFromUI() (string, string, todo.NullTime, error) {
	taskNamePrompt := promptui.Prompt{
		Label: "Task Name",
	}
	taskName, err := taskNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", "", todo.NullTime{}, err
	}

	projectNamePrompt := promptui.Prompt{
		Label: "Project Name",
	}
	projectName, err := projectNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", "", todo.NullTime{}, err
	}

	duePrompt := promptui.Prompt{
		Label: "Due Date (YYYY-MM-DD)",
	}
	dueString, err := duePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", "", todo.NullTime{}, err
	}
	due, err := parseNullTimeString(dueString)
	if err != nil {
		return "", "", todo.NullTime{}, fmt.Errorf("Error parsing time: %w\n", err)
	}

	return taskName, projectName, due, nil
}

func parseNullTimeString(timeString string) (todo.NullTime, error) {
	if timeString == "" {
		return todo.NullTime{}, nil
	}

	t, err := time.ParseInLocation("2006-01-02", timeString, time.Local)
	if err != nil {
		return todo.NullTime{}, fmt.Errorf("Error parsing time: %w\n", err)
	}

	return todo.NullTime{
		Time: t,
		Valid: true,
	}, nil
}