package grpc

import (
	"context"

	"github.com/reiot777/spansqlx-example/grpc/packet"
	"github.com/reiot777/spansqlx-example/store"
)

var _ packet.TodoServiceServer = (*Service)(nil)

func (s *Service) CreateTask(ctx context.Context, in *packet.CreateTaskRequest) (*packet.CreateTaskResponse, error) {
	newTask := &store.Task{
		AccountID: in.AccountId,
		Name:      in.Name,
	}
	if err := s.Store.Tasks().AddTask(ctx, newTask); err != nil {
		return nil, err
	}
	return packet.NewCreateTaskResponse(newTask), nil
}

func (s *Service) ListTasks(ctx context.Context, in *packet.ListTasksRequest) (*packet.ListTasksResponse, error) {
	tasks, err := s.Store.Tasks().Tasks(ctx, store.TaskQuery{
		AccountID: &in.AccountId,
	})
	if err != nil {
		return nil, err
	}
	return packet.NewListTasksResponse(tasks), nil
}
