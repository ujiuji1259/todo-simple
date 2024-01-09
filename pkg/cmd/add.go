package cmd

import (
	"fmt"

	"todo-simple/pkg/manager"
	"todo-simple/pkg/db"
	"github.com/manifoldco/promptui"
	"github.com/adrg/xdg"
)

type AddCmd struct {}

func (c *AddCmd) Run() error {
	taskName, projectName, err := loadTodoFromUI()
	if err != nil {
		fmt.Printf("Error loading todo from UI: %v\n", err)
		return err
	}

	todoFilePath, err := xdg.DataFile("todo-simple/todo.tsv")
	if err != nil {
		return err
	}
	todoFile, err := db.NewTodoFile(todoFilePath)
	if err != nil {
		return fmt.Errorf("Error loading todo file: %v", err)
	}
	todoManager, err := manager.NewTodoManager(todoFile)
	if err != nil {
		fmt.Printf("Error loading todo manager: %v\n", err)
		return err
	}

	err = todoManager.Add(taskName, projectName)
	if err != nil {
		fmt.Printf("Error loading todo manager: %v\n", err)
		return err
	}

	err = todoManager.Save()
	if err != nil {
		fmt.Printf("Error saving todo manager: %v\n", err)
		return err
	}
	return nil
}

func loadTodoFromUI() (string, string, error) {
	taskNamePrompt := promptui.Prompt{
		Label:    "Task Name",
	}
	taskName, err := taskNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", "", err
	}

	projectNamePrompt := promptui.Prompt{
		Label:    "Project Name",
	}
	projectName, err := projectNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", "", err
	}

	return taskName, projectName, nil
}