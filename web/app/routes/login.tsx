import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/login')({
  beforeLoad: ({ context }) => {
    if (context.auth.isAuthenticated) {
      throw redirect({ to: '/todos' })
    }
  },
  component: LoginPage,
})

function LoginPage() {
  const apiUrl = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
  return (
    <div>
      <h1>Todo App</h1>
      <a href={`${apiUrl}/auth/login`}>
        <button type="button">Sign in with Google</button>
      </a>
    </div>
  )
}
