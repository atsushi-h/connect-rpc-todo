import { createClient } from '@connectrpc/connect'
import { TransportProvider } from '@connectrpc/connect-query'
import { createGrpcWebTransport } from '@connectrpc/connect-web'
import { type QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createRootRouteWithContext, HeadContent, Outlet, Scripts } from '@tanstack/react-router'
import { createIsomorphicFn } from '@tanstack/react-start'
import { AuthService } from '@todo-app/api-client/src/auth/v1/auth_pb.js'
import type { ReactNode } from 'react'
import { env } from '~/lib/env'
import { transport } from '~/lib/transport'

interface RouterContext {
  queryClient: QueryClient
}

export interface AuthContext {
  isAuthenticated: boolean
  user?: {
    id: string
    email: string
    displayName: string
    avatarUrl: string
  }
}

const API_BASE = env.VITE_API_URL

// SSR 時はリクエストの Cookie をバックエンドへ転送するトランスポートを使用する
const getAuthTransport = createIsomorphicFn()
  .client(() => transport)
  .server(async () => {
    const { getRequestHeader } = await import('@tanstack/react-start/server')
    const cookieHeader = getRequestHeader('cookie') ?? ''
    return createGrpcWebTransport({
      baseUrl: API_BASE,
      fetch: (input, init) => {
        const headers = new Headers(init?.headers)
        if (cookieHeader) headers.set('cookie', cookieHeader)
        return fetch(input as string, { ...(init as RequestInit), headers })
      },
    })
  })

function NotFound() {
  return <p>Not Found</p>
}

export const Route = createRootRouteWithContext<RouterContext>()({
  notFoundComponent: NotFound,
  beforeLoad: async (): Promise<{ auth: AuthContext }> => {
    const authTransport = await getAuthTransport()
    const client = createClient(AuthService, authTransport)
    try {
      const me = await client.getMe({})
      return {
        auth: {
          isAuthenticated: true,
          user: {
            id: me.id,
            email: me.email,
            displayName: me.displayName,
            avatarUrl: me.avatarUrl,
          },
        },
      }
    } catch {
      return { auth: { isAuthenticated: false } }
    }
  },
  component: RootComponent,
})

function RootDocument({ children }: { children: ReactNode }) {
  return (
    <html lang="ja">
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <Scripts />
      </body>
    </html>
  )
}

function RootComponent() {
  const { queryClient } = Route.useRouteContext()
  return (
    <RootDocument>
      <TransportProvider transport={transport}>
        <QueryClientProvider client={queryClient}>
          <Outlet />
        </QueryClientProvider>
      </TransportProvider>
    </RootDocument>
  )
}
