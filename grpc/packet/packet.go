package packet

import (
	"time"

	"github.com/reiot777/spansqlx-example/store"
)

func NewCreateAccountResponse(arg *store.Account) *CreateAccountResponse {
	return &CreateAccountResponse{
		Data: NewAccount(arg),
	}
}

func NewGetAccountResponse(arg *store.Account) *GetAccountResponse {
	return &GetAccountResponse{
		Data: NewAccount(arg),
	}
}

func NewListAccountsResponse(args []*store.Account) *ListAccountsResponse {
	data := make([]*Account, len(args))
	for i := range args {
		data[i] = NewAccount(args[i])
	}
	return &ListAccountsResponse{
		Data: data,
	}
}

func NewCreateTaskResponse(arg *store.Task) *CreateTaskResponse {
	return &CreateTaskResponse{
		Data: NewTask(arg),
	}
}

func NewListTasksResponse(args []*store.Task) *ListTasksResponse {
	data := make([]*Task, len(args))
	for i := range args {
		data[i] = NewTask(args[i])
	}
	return &ListTasksResponse{
		Data: data,
	}
}

func NewAccount(arg *store.Account) *Account {
	return &Account{
		Id:        arg.ID,
		Owner:     arg.Owner,
		Email:     arg.Email,
		CreatedAt: arg.CreatedAt.Format(time.RFC3339),
		UpdatedAt: arg.UpdatedAt.Format(time.RFC3339),
	}
}

func NewTask(arg *store.Task) *Task {
	return &Task{
		Id:        arg.ID,
		Name:      arg.Name,
		Undone:    arg.Undone,
		CreatedAt: arg.CreatedAt.Format(time.RFC3339),
		UpdatedAt: arg.UpdatedAt.Format(time.RFC3339),
	}
}
