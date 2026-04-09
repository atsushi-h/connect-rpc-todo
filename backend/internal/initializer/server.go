package initializer

import (
	"context"
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"gen/go/todo/v1/todov1connect"
	"todo-app/backend/internal/config"
	"todo-app/backend/internal/db"
	"todo-app/backend/internal/handler"
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
	todoHandler := handler.NewTodoHandler(queries)

	mux := http.NewServeMux()
	path, h := todov1connect.NewTodoServiceHandler(todoHandler)
	mux.Handle(path, h)

	srv := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	return srv, cfg, cleanup, nil
}
