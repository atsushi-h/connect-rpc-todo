import { createClient } from '@connectrpc/connect'
import { AuthService } from '@todo-app/api-client/src/auth/v1/auth_pb.js'
import * as AuthSession from 'expo-auth-session'
import { router } from 'expo-router'
import { useEffect } from 'react'
import { Button, Text, View } from 'react-native'
import { env } from '../../lib/env'
import { storage } from '../../lib/storage'
import { transport } from '../../lib/transport'

// iOS OAuth クライアント用 reverse client ID スキーム
// 例: com.googleusercontent.apps.338570704218-xxx:/
// Google は iOS OAuth クライアントに対してこの形式を Console 登録なしで自動許可する
const clientId = env.EXPO_PUBLIC_GOOGLE_CLIENT_ID
const reverseClientId = `com.googleusercontent.apps.${clientId.replace('.apps.googleusercontent.com', '')}`
// native パラメータで明示的に指定（expo-auth-session Google プロバイダと同じパターン）
const redirectUri = AuthSession.makeRedirectUri({ native: `${reverseClientId}:/oauthredirect` })
console.log('[OAuth] redirectUri:', redirectUri)

export default function LoginScreen() {
  // Google OAuth のエンドポイントを自動検出（Hook はコンポーネント内で呼ぶ）
  const discovery = AuthSession.useAutoDiscovery('https://accounts.google.com')

  const [request, response, promptAsync] = AuthSession.useAuthRequest(
    {
      clientId: env.EXPO_PUBLIC_GOOGLE_CLIENT_ID,
      scopes: ['openid', 'profile', 'email'],
      codeChallengeMethod: AuthSession.CodeChallengeMethod.S256, // PKCE
      redirectUri,
    },
    discovery,
  )

  useEffect(() => {
    if (response?.type !== 'success') return

    const { code } = response.params
    const codeVerifier = request?.codeVerifier // expo-auth-session が生成

    if (!code || !codeVerifier) return

    // ExchangeToken RPC で JWT 取得
    const client = createClient(AuthService, transport)
    client
      .exchangeToken({ code, codeVerifier, redirectUri })
      .then(async ({ accessToken }) => {
        await storage.setToken(accessToken)
        router.replace('/(app)/todos')
      })
      .catch(console.error)
  }, [response, request])

  return (
    <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
      <Text style={{ fontSize: 24, fontWeight: 'bold', marginBottom: 24 }}>Todo App</Text>
      <Button title="Sign in with Google" onPress={() => promptAsync()} disabled={!request} />
    </View>
  )
}
