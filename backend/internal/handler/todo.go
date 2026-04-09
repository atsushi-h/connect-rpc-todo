package handler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	todov1 "gen/go/todo/v1"
	"gen/go/todo/v1/todov1connect"
	"todo-app/backend/internal/db"
)

type TodoHandler struct {
	todov1connect.UnimplementedTodoServiceHandler
	queries *db.Queries
}

func NewTodoHandler(queries *db.Queries) *TodoHandler {
	return &TodoHandler{queries: queries}
}

// TODO(phase4): context.Value(middleware.UserIDKey) に差し替える
var tempUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func todoToProto(t db.Todo) *todov1.Todo {
	return &todov1.Todo{
		Id:        t.ID.String(),
		Title:     t.Title,
		Completed: t.Completed,
	}
}

func (h *TodoHandler) ListTodos(
	ctx context.Context,
	req *connect.Request[todov1.ListTodosRequest],
) (*connect.Response[todov1.ListTodosResponse], error) {
	todos, err := h.queries.ListTodosByUser(ctx, tempUserID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("list todos: %w", err))
	}
	protoTodos := make([]*todov1.Todo, len(todos))
	for i, t := range todos {
		protoTodos[i] = todoToProto(t)
	}
	return connect.NewResponse(&todov1.ListTodosResponse{Todos: protoTodos}), nil
}

func (h *TodoHandler) CreateTodo(
	ctx context.Context,
	req *connect.Request[todov1.CreateTodoRequest],
) (*connect.Response[todov1.CreateTodoResponse], error) {
	todo, err := h.queries.CreateTodo(ctx, db.CreateTodoParams{
		UserID: tempUserID,
		Title:  req.Msg.Title,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("create todo: %w", err))
	}
	return connect.NewResponse(&todov1.CreateTodoResponse{Todo: todoToProto(todo)}), nil
}

func (h *TodoHandler) UpdateTodo(
	ctx context.Context,
	req *connect.Request[todov1.UpdateTodoRequest],
) (*connect.Response[todov1.UpdateTodoResponse], error) {
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid id: %w", err))
	}
	todo, err := h.queries.UpdateTodo(ctx, db.UpdateTodoParams{
		ID:        id,
		UserID:    tempUserID,
		Title:     req.Msg.Title,
		Completed: req.Msg.Completed,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("update todo: %w", err))
	}
	return connect.NewResponse(&todov1.UpdateTodoResponse{Todo: todoToProto(todo)}), nil
}

func (h *TodoHandler) DeleteTodo(
	ctx context.Context,
	req *connect.Request[todov1.DeleteTodoRequest],
) (*connect.Response[todov1.DeleteTodoResponse], error) {
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid id: %w", err))
	}
	if err := h.queries.DeleteTodo(ctx, db.DeleteTodoParams{
		ID:     id,
		UserID: tempUserID,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("delete todo: %w", err))
	}
	return connect.NewResponse(&todov1.DeleteTodoResponse{}), nil
}
