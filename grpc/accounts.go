package grpc

import (
	"context"
	"errors"
	"strings"

	"github.com/reiot777/spansqlx-example/grpc/packet"
	"github.com/reiot777/spansqlx-example/store"
)

var _ packet.AccountServiceServer = (*Service)(nil)

func (s *Service) CreateAccount(ctx context.Context, in *packet.CreateAccountRequest) (*packet.CreateAccountResponse, error) {
	newAccount := &store.Account{
		Owner: strings.Split(in.Email, "@")[0],
		Email: in.Email,
	}

	if err := s.Store.Pipeline(ctx, func(ctx context.Context) error {
		_, err := s.Store.Accounts().Account(ctx, store.AccountQuery{})
		if err != nil && !errors.Is(err, store.ErrEmailAlreadyExists) {
			return err
		}

		if err := s.Store.Accounts().AddAccount(ctx, newAccount); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return packet.NewCreateAccountResponse(newAccount), nil
}

func (s *Service) GetAccount(ctx context.Context, in *packet.GetAccountRequest) (*packet.GetAccountResponse, error) {
	account, err := s.Store.Accounts().Account(ctx, store.AccountQuery{
		Owner: &in.Owner,
	})
	if err != nil {
		return nil, err
	}

	return packet.NewGetAccountResponse(account), nil
}

func (s *Service) ListAccounts(ctx context.Context, in *packet.ListAccountsRequest) (*packet.ListAccountsResponse, error) {
	accounts, err := s.Store.Accounts().Accounts(ctx, store.AccountQuery{
		Range: &store.Range{
			Limit:  in.PageSize,
			Offset: in.NextPage,
		},
	})
	if err != nil {
		return nil, err
	}
	return packet.NewListAccountsResponse(accounts), nil
}
