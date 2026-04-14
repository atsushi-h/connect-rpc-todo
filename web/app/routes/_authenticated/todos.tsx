import { createFileRoute } from "@tanstack/react-router";
import { useQuery, useMutation } from "@connectrpc/connect-query";
import { useQueryClient } from "@tanstack/react-query";
import { TodoService } from "@todo-app/api-client/src/todo/v1/todo_pb.js";
import { createClient } from "@connectrpc/connect";
import { AuthService } from "@todo-app/api-client/src/auth/v1/auth_pb.js";
import { transport } from "~/lib/transport";
import { useState } from "react";

export const Route = createFileRoute("/_authenticated/todos")({
  component: TodosPage,
});

function TodosPage() {
  const { auth } = Route.useRouteContext();
  const queryClient = useQueryClient();
  const [newTitle, setNewTitle] = useState("");

  const { data, isLoading, error } = useQuery(TodoService.method.listTodos, {});

  const createMutation = useMutation(TodoService.method.createTodo, {
    onSuccess: () => {
      queryClient.invalidateQueries();
      setNewTitle("");
    },
  });

  const updateMutation = useMutation(TodoService.method.updateTodo, {
    onSuccess: () => queryClient.invalidateQueries(),
  });

  const deleteMutation = useMutation(TodoService.method.deleteTodo, {
    onSuccess: () => queryClient.invalidateQueries(),
  });

  const handleSignOut = async () => {
    const client = createClient(AuthService, transport);
    await client.signOut({});
    window.location.href = "/login";
  };

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <header>
        <h1>Todos</h1>
        <span>{auth.user?.displayName}</span>
        <button onClick={handleSignOut}>Sign Out</button>
      </header>

      <form
        onSubmit={(e) => {
          e.preventDefault();
          if (newTitle.trim()) createMutation.mutate({ title: newTitle });
        }}
      >
        <input
          value={newTitle}
          onChange={(e) => setNewTitle(e.target.value)}
          placeholder="New todo..."
        />
        <button type="submit" disabled={createMutation.isPending}>
          Add
        </button>
      </form>

      <ul>
        {data?.todos.map((todo) => (
          <li key={todo.id}>
            <input
              type="checkbox"
              checked={todo.completed}
              onChange={() =>
                updateMutation.mutate({
                  id: todo.id,
                  title: todo.title,
                  completed: !todo.completed,
                })
              }
            />
            <span
              style={{ textDecoration: todo.completed ? "line-through" : "none" }}
            >
              {todo.title}
            </span>
            <button onClick={() => deleteMutation.mutate({ id: todo.id })}>
              Delete
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
