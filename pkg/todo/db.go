package todo

import (
	"context"
	"fmt"

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

func NewTsvDb(path string) (*TsvDb, error) {
	db, err := sqlx.Open("csvq", path)
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
	sql := sq.Select("id", "task", "project", "status", "due").From("`todo.tsv`")
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
		Columns("id", "task", "project", "status", "due").
		Values(todoItem.Id, todoItem.TaskName, todoItem.Project, todoItem.Status.String(), todoItem.Due).
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
	query, args, err := sq.Select("id", "task", "project", "status", "due").
		From("`todo.tsv`").
		Where(sq.Eq{"id": taskId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	listQueryWithDelemiter := addDelemiterToQuery(query)

	var item TodoItem
	err = db.Db.GetContext(ctx, &item, listQueryWithDelemiter, args...)
	if err != nil {
		return nil, fmt.Errorf("fail to load: %w", err)
	}

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
