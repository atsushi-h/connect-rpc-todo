import type { Interceptor } from '@connectrpc/connect'
import { createConnectTransport } from '@connectrpc/connect-web'
import { Platform } from 'react-native'
import { storage } from './storage'

const getBaseUrl = () => {
  // Android Emulator は仮想 NIC 経由で 10.0.2.2 がホスト Mac の localhost
  if (Platform.OS === 'android') return 'http://10.0.2.2:8080'
  return 'http://localhost:8080'
}

const authInterceptor: Interceptor = (next) => async (req) => {
  const token = await storage.getToken()
  if (token) {
    req.header.set('authorization', `Bearer ${token}`)
  }
  return next(req)
}

export const transport = createConnectTransport({
  baseUrl: getBaseUrl(),
  interceptors: [authInterceptor],
})
