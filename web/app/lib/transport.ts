import { createConnectTransport } from '@connectrpc/connect-web'
import { env } from '~/lib/env'

export const transport = createConnectTransport({
  baseUrl: env.VITE_API_URL,
  // HttpOnly Cookie を自動送信するために credentials: "include" を fetch でオーバーライド
  fetch: (input, init) => fetch(input, { ...init, credentials: 'include' }),
})
