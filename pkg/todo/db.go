package todo

import (
	"context"
	"database/sql"
	"fmt"

	// "github.com/mithrandie/csvq/lib/query"

	_ "github.com/mithrandie/csvq-driver"
)

type Db interface {
	List() []*TodoItem
}

type TsvDb struct {
	Db *sql.DB
}

func NewTsvDb(ctx context.Context, path string) (*TsvDb, error) {
	db, err := sql.Open("csvq", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open csvq: %w", err)
	}

	return &TsvDb{
		Db: db,
	}, nil
}

func (db *TsvDb) List(ctx context.Context) ([]*TodoItem, error) {
	rows, err := db.Db.QueryContext(ctx, "SET @@DELIMITER TO '\t'; SELECT * FROM `todo.tsv`")
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var items []*TodoItem
	for rows.Next() {
		var item TodoItem
		var status string
		if err := rows.Scan(&item.Id, &item.TaskName, &item.Project, &status); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		todoStatus, err := TodoStatusString(status)
		if err != nil {
			return nil, fmt.Errorf("invalid todo status: %w", err)
		}
		item.Status = todoStatus

		items = append(items, &item)
	}

	return items, nil
}