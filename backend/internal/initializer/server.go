package initializer

import (
	"context"
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"connectrpc.com/connect"
	"gen/go/auth/v1/authv1connect"
	"gen/go/todo/v1/todov1connect"
	"todo-app/backend/internal/config"
	"todo-app/backend/internal/db"
	"todo-app/backend/internal/handler"
	"todo-app/backend/internal/middleware"
)

func BuildServer(_ context.Context) (*http.Server, *config.Config, func(), error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, func() {}, err
	}

	conn, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, nil, func() {}, err
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, nil, func() {}, err
	}
	cleanup := func() {
		conn.Close()
	}

	queries := db.New(conn)

	interceptors := connect.WithInterceptors(middleware.NewAuthInterceptor(cfg))

	mux := http.NewServeMux()

	// OAuth コールバック（Web フロー）
	authCBHandler := handler.NewAuthCallbackHandler(queries, cfg)
	mux.Handle("GET /auth/login", authCBHandler.Login())
	mux.Handle("GET /auth/callback", authCBHandler.Callback())

	// AuthService RPC（Native フロー + GetMe/SignOut）
	authRPCHandler := handler.NewAuthHandler(queries, cfg)
	path2, h2 := authv1connect.NewAuthServiceHandler(authRPCHandler, interceptors)
	mux.Handle(path2, h2)

	// TodoService（interceptor 付き）
	todoHandler := handler.NewTodoHandler(queries)
	path, h := todov1connect.NewTodoServiceHandler(todoHandler, interceptors)
	mux.Handle(path, h)

	srv := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	return srv, cfg, cleanup, nil
}
