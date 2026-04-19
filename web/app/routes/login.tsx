import { createFileRoute, redirect } from '@tanstack/react-router'
import { env } from '~/lib/env'

export const Route = createFileRoute('/login')({
  beforeLoad: ({ context }) => {
    if (context.auth.isAuthenticated) {
      throw redirect({ to: '/todos' })
    }
  },
  component: LoginPage,
})

function LoginPage() {
  const apiUrl = env.VITE_API_URL
  return (
    <div>
      <h1>Todo App</h1>
      <a href={`${apiUrl}/auth/login`}>
        <button type="button">Sign in with Google</button>
      </a>
    </div>
  )
}
