package cmd

import (
	"context"
	"fmt"
	"time"
)

type StartCmd struct {
	Id     string          `help:"Task Id"`
}

func (c *StartCmd) Run() error {
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

	err = db.Start(ctx, c.Id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to start %s: %w", c.Id, err)
	}
	return nil
}
