package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	sqlr "github.com/Masterminds/squirrel"
	"github.com/reiot777/spansqlx"
	"github.com/rs/xid"
)

var _ AccountStore = (*accountStore)(nil)

type accountStore struct {
	*store
}

func (s *accountStore) AddAccount(ctx context.Context, arg *Account) error {
	if arg == nil {
		return fmt.Errorf("Account provide is nil")
	}

	if arg.Email == "" || !strings.Contains(arg.Email, "@") {
		return fmt.Errorf("Account.Email not be nil")
	}

	if arg.ID == "" {
		arg.ID = xid.New().String()
	}
	if arg.Owner == "" {
		arg.Owner = strings.Split(arg.Email, "@")[0]
	}
	if arg.CreatedAt.IsZero() {
		arg.CreatedAt = time.Now()
	}
	if arg.UpdatedAt.IsZero() {
		arg.UpdatedAt = arg.CreatedAt
	}

	return s.db.Exec(ctx,
		`
		INSERT INTO accounts (id, owner, email, createdAt, updatedAt)
		VALUES (@id, @owner, @email, @createdAt, @updatedAt)
		`,
		arg.ID,
		arg.Owner,
		arg.Email,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
}

func (s *accountStore) Account(ctx context.Context, q AccountQuery) (*Account, error) {
	var account Account

	sql, args := s.genQuerySql(q)

	if err := s.db.Get(ctx, &account, sql, args...); err != nil {
		if errors.Is(err, spansqlx.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return &account, nil
}

func (s *accountStore) Accounts(ctx context.Context, q AccountQuery) ([]*Account, error) {
	var accounts []*Account
	sql, args := s.genQuerySql(q)

	if err := s.db.Select(ctx, &accounts, sql, args...); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *accountStore) genQuerySql(q AccountQuery) (sql string, args []interface{}) {
	var (
		table = sqlr.Select("*").From("accounts")
		where []sqlr.Eq
	)

	if v := q.Owner; v != nil {
		where = append(where, sqlr.Eq{"Owner": *v})
	}
	if v := q.Email; v != nil {
		where = append(where, sqlr.Eq{"Email": *v})
	}

	// apply where conditions
	for i := range where {
		table = table.Where(where[i])
	}

	// apply limit and offset
	if v := q.Range; v != nil {
		table = table.Limit(uint64(v.Limit)).Offset((uint64(v.Offset) - 1) * uint64(v.Limit))
	}

	sql, args, _ = table.PlaceholderFormat(sqlr.AtP).ToSql()
	return
}
