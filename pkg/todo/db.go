package todo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"encoding/csv"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	sq "github.com/Masterminds/squirrel"

	_ "github.com/mithrandie/csvq-driver"
)

type Db interface {
	List() []*TodoItem
}

type TsvDb struct {
	Db *sqlx.DB
}

func initTodoFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	headers := []string{"id", "task", "project", "status", "due", "estimation", "started_at", "ended_at"}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open todo.tsv: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	defer writer.Flush()
	err = writer.Write(headers)
	if err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	return nil
}

func NewTsvDb(path string) (*TsvDb, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create todo-simple directory: %w", err)
		}
	}

	tsvFiles := []string{filepath.Join(path, "todo.tsv")}
	for _, tsvFile := range tsvFiles {
		err := initTodoFile(tsvFile)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize %s: %w", tsvFile, err)
		}
	}

	db, err := sqlx.Open("csvq", path)
	if err != nil {
		return nil, fmt.Errorf("failed to opening tsv: %w", err)
	}

	return &TsvDb{Db: db}, nil
}

func (db *TsvDb) Close() error {
	return db.Db.Close()
}

func addDelemiterToQuery(query string) string {
	return "SET @@DELIMITER TO '\t'; " + query
}

func (db *TsvDb) ListItems(ctx context.Context, projects []string, statuses []TodoStatus) ([]*TodoItem, error) {
	sql := sq.Select("*").From("`todo.tsv`")
	if len(projects) > 0 {
		sql = sql.Where(sq.Eq{"project": projects})
	}
	if len(statuses) > 0 {
		sql = sql.Where(sq.Eq{"status": StatusStrings(statuses)})
	}

	query, args, err := sql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var items []*TodoItem
	listQueryWithDelemiter := addDelemiterToQuery(query)
	err = db.Db.SelectContext(ctx, &items, listQueryWithDelemiter, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	return items, nil
}

func (db *TsvDb) Add(ctx context.Context, todoItem TodoItem) error {
	query, args, err := sq.Insert("`todo.tsv`").
		Columns("id", "task", "project", "status", "due", "estimation", "started_at", "ended_at").
		Values(todoItem.Id, todoItem.TaskName, todoItem.Project, todoItem.Status, todoItem.Due, todoItem.Estimation, todoItem.StartedAt, todoItem.EndedAt).
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
	if affected == 0 {
		fmt.Printf("RowsAffected: %d\n", affected)
		return errors.New("target task is not found")
	}
	fmt.Printf("RowsAffected: %d\n", affected)

	return nil
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
