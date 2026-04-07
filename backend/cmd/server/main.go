package main

import (
	"fmt"
	"gen/go/todo/v1/todov1connect"
	"log"
	"net/http"
	"todo-app/backend/internal/handler"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()
	todoHandler := &handler.TodoHandler{}
	path, h := todov1connect.NewTodoServiceHandler(todoHandler)
	mux.Handle(path, h)
	addr := ":8080"
	fmt.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	))
}
