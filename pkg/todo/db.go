package todo

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/xid"

	_ "github.com/mithrandie/csvq-driver"
)

type Db interface {
	List() []*TodoItem
}

type TsvDb struct {
	Db *sql.DB
}

func NewTsvDb(path string) (*TsvDb, error) {
	db, err := sql.Open("csvq", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open csvq: %w", err)
	}

	return &TsvDb{
		Db: db,
	}, nil
}

func (db *TsvDb) Close() error {
	return db.Db.Close()
}

func addDelemiterToQuery(query string) string {
	return "SET @@DELIMITER TO '\t'; " + query
}

func (db *TsvDb) ListItems(ctx context.Context, projects []string, statuses []TodoStatus) ([]*TodoItem, error) {
	sql := sq.Select("id", "task", "project", "status").From("`todo.tsv`")
	if len(projects) > 0 {
		sql = sql.Where(sq.Eq{"project": projects})
	}
	if len(statuses) > 0 {
		sql = sql.Where(sq.Eq{"status": statusStrings(statuses)})
	}

	query, args, err := sql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	listQueryWithDelemiter := addDelemiterToQuery(query)
	rows, err := db.Db.QueryContext(ctx, listQueryWithDelemiter, args...)
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

func (db *TsvDb) Add(ctx context.Context, taskName string, projectName string) error {
	guid := xid.New()

	todoStatus, err := TodoStatusString("Todo")
	if err != nil {
		return fmt.Errorf("Error initializing todo status: %w\n", err)
	}
	query, args, err := sq.Insert("`todo.tsv`").
		Columns("id", "task", "project", "status").
		Values(guid.String(), taskName, projectName, todoStatus.String()).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	ret, err := db.Db.ExecContext(ctx, addDelemiterToQuery(query), args...)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	affected, _ := ret.RowsAffected()
	fmt.Printf("RowsAffected: %d\n", affected)

	return nil
}

func (db *TsvDb) Delete(ctx context.Context, taskId string) error {
	query, args, err := sq.Delete("`todo.tsv`").
		Where(sq.Eq{"id": taskId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	ret, err := db.Db.ExecContext(ctx, addDelemiterToQuery(query), args...)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	affected, _ := ret.RowsAffected()
	fmt.Printf("RowsAffected: %d\n", affected)

	return nil
}

func (db *TsvDb) ListItem(ctx context.Context, taskId string) (*TodoItem, error) {
	query, args, err := sq.Select("id", "task", "project", "status").
		From("`todo.tsv`").
		Where(sq.Eq{"id": taskId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	listQueryWithDelemiter := addDelemiterToQuery(query)
	row := db.Db.QueryRowContext(ctx, listQueryWithDelemiter, args...)

	var item TodoItem
	var status string
	if err := row.Scan(&item.Id, &item.TaskName, &item.Project, &status); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}
	todoStatus, err := TodoStatusString(status)
	if err != nil {
		return nil, fmt.Errorf("invalid todo status: %w", err)
	}
	item.Status = todoStatus

	return &item, nil
}

func (db *TsvDb) UpdateStatus(ctx context.Context, taskId string, status TodoStatus) error {
	query, args, err := sq.Update("`todo.tsv`").
		Set("status", status.String()).
		Where(sq.Eq{"id": taskId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	ret, err := db.Db.ExecContext(ctx, addDelemiterToQuery(query), args...)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	affected, _ := ret.RowsAffected()
	fmt.Printf("RowsAffected: %d\n", affected)

	return nil
}
