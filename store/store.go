package store

import (
	"context"
	"time"

	"github.com/reiot777/spansqlx"
)

type Range struct {
	Limit  int32
	Offset int32
}

type Account struct {
	ID        string
	Owner     string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountQuery struct {
	Owner *string
	Email *string
	Range *Range
}

type AccountStore interface {
	AddAccount(context.Context, *Account) error
	Account(context.Context, AccountQuery) (*Account, error)
	Accounts(context.Context, AccountQuery) ([]*Account, error)
}

type Task struct {
	ID        string
	AccountID string
	Name      string
	Undone    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TaskQuery struct {
	AccountID *string
	Undone    *bool
}

type TaskStore interface {
	Task(context.Context, TaskQuery) (*Task, error)
	Tasks(context.Context, TaskQuery) ([]*Task, error)
	AddTask(context.Context, *Task) error
	SetTaskDone(context.Context, string) error
}

type Store interface {
	Pipeline(context.Context, func(context.Context) error) error
	Accounts() AccountStore
	Tasks() TaskStore
}

var _ Store = (*store)(nil)

type store struct {
	db *spansqlx.DB
}

func New(db *spansqlx.DB) Store {
	return &store{db: db}
}

func (s *store) Pipeline(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.db.TxPipeline(ctx, fn)
}

func (s *store) Accounts() AccountStore {
	return &accountStore{}
}

func (s *store) Tasks() TaskStore {
	return &taskStore{}
}
