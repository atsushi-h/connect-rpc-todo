package handler

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	todov1 "gen/go/todo/v1"
	"gen/go/todo/v1/todov1connect"
	"todo-app/backend/internal/db"
	"todo-app/backend/internal/middleware"
)

type TodoHandler struct {
	todov1connect.UnimplementedTodoServiceHandler
	queries *db.Queries
}

func NewTodoHandler(queries *db.Queries) *TodoHandler {
	return &TodoHandler{queries: queries}
}

func userIDFromCtx(ctx context.Context) (uuid.UUID, error) {
	str, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return uuid.UUID{}, errors.New("unauthenticated")
	}
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user id: %w", err)
	}
	return id, nil
}

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
	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	todos, err := h.queries.ListTodosByUser(ctx, userID)
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
	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	todo, err := h.queries.CreateTodo(ctx, db.CreateTodoParams{
		UserID: userID,
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
	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid id: %w", err))
	}
	todo, err := h.queries.UpdateTodo(ctx, db.UpdateTodoParams{
		ID:        id,
		UserID:    userID,
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
	userID, err := userIDFromCtx(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid id: %w", err))
	}
	if err := h.queries.DeleteTodo(ctx, db.DeleteTodoParams{
		ID:     id,
		UserID: userID,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("delete todo: %w", err))
	}
	return connect.NewResponse(&todov1.DeleteTodoResponse{}), nil
}
