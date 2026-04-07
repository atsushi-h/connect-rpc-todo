package handler

import (
	"context"
	todov1 "gen/go/todo/v1"
	"gen/go/todo/v1/todov1connect"

	"connectrpc.com/connect"
)

type TodoHandler struct {
	todov1connect.UnimplementedTodoServiceHandler
}

func (h *TodoHandler) ListTodos(
	ctx context.Context,
	req *connect.Request[todov1.ListTodosRequest],
) (*connect.Response[todov1.ListTodosResponse], error) {
	return connect.NewResponse(&todov1.ListTodosResponse{
		Todos: []*todov1.Todo{
			{Id: "1", Title: "Buy milk", Completed: false},
			{Id: "2", Title: "Write code", Completed: true},
		},
	}), nil
}
