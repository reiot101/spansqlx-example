package store

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/spanner"
	sqlr "github.com/Masterminds/squirrel"
	"github.com/reiot777/spansqlx"
	"github.com/rs/xid"
)

var _ TaskStore = (*taskStore)(nil)

type taskStore struct {
	*store
}

func (s *taskStore) Task(ctx context.Context, q TaskQuery) (*Task, error) {
	var task Task

	sql, args := s.genQuerySql(q)

	if err := s.db.Get(ctx, &task, sql, args...); err != nil {
		if errors.Is(err, spansqlx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (s *taskStore) Tasks(ctx context.Context, q TaskQuery) ([]*Task, error) {
	var tasks []*Task
	sql, args := s.genQuerySql(q)

	if err := s.db.Select(ctx, &tasks, sql, args...); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *taskStore) AddTask(ctx context.Context, task *Task) error {
	if task.ID == "" {
		task.ID = xid.New().String()
	}
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = task.CreatedAt
	}

	return s.db.NamedExec(ctx,
		`
		INSERT INTO Tasks (
			ID, 
			Name, 
			Undone, 
			AccountID, 
			CreatedAt, 
			UpdatedAt
		) 
		VALUES (
			@ID, 
			@Name, 
			@Undone, 
			@AccountID, 
			@CreatedAt, 
			@UpdatedAt
		)
		`,
		task)
}

func (s *taskStore) SetTaskDone(ctx context.Context, id string) error {
	var n int64
	if err := s.db.Get(ctx, &n, `SELECT COUNT(1) FROM Tasks WHERE ID=@ID`, id); err != nil {
		return err
	}
	if n == 0 {
		return ErrTaskNotFound
	}

	return s.db.ExecX(ctx, spanner.Statement{
		SQL: "UPDATE Tasks SET Undone=@Undone, UpdatedAt=@UpdatedAt WHERE ID=@ID",
		Params: map[string]interface{}{
			"Undone":    false,
			"UpdatedAt": time.Now(),
			"ID":        id,
		},
	})
}

func (s *taskStore) genQuerySql(q TaskQuery) (sql string, args []interface{}) {
	var (
		table = sqlr.Select("*").From("tasks")
		where []sqlr.Eq
	)

	if v := q.AccountID; v != nil {
		where = append(where, sqlr.Eq{"AccountID": *v})
	}
	if v := q.Undone; v != nil {
		where = append(where, sqlr.Eq{"Undone": *v})
	}

	// apply where conditions
	for i := range where {
		table = table.Where(where[i])
	}

	sql, args, _ = table.PlaceholderFormat(sqlr.AtP).ToSql()
	return
}
