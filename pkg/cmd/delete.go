package cmd

import (
	"context"
	"fmt"
	"time"
)

type DeleteCmd struct {
	Id string `help:"Id"`
}

func (c *DeleteCmd) Run() error {
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

	err = db.Delete(ctx, c.Id)
	if err != nil {
		return fmt.Errorf("Error deleting %s: %v\n", c.Id, err)
	}

	return nil
}
