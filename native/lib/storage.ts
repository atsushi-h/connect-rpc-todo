import * as SecureStore from 'expo-secure-store'

const JWT_KEY = 'jwt_token'

export const storage = {
  getToken: () => SecureStore.getItemAsync(JWT_KEY),
  setToken: (token: string) => SecureStore.setItemAsync(JWT_KEY, token),
  deleteToken: () => SecureStore.deleteItemAsync(JWT_KEY),
}
