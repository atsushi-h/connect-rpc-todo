// app.json の代わりに動的設定ファイルを使用（env でスキームを制御）
const clientId = process.env.EXPO_PUBLIC_GOOGLE_CLIENT_ID ?? ''
// "338570704218-xxx.apps.googleusercontent.com" → "com.googleusercontent.apps.338570704218-xxx"
const reverseClientId = clientId
  ? `com.googleusercontent.apps.${clientId.replace('.apps.googleusercontent.com', '')}`
  : ''

/** @type {import('expo/config').ExpoConfig} */
module.exports = {
  name: 'Todo App',
  slug: 'todo-app',
  scheme: ['todoapp', reverseClientId].filter(Boolean),
  version: '1.0.0',
  ios: {
    bundleIdentifier: 'com.example.todoapp',
    infoPlist: {
      // iOS が reverse client ID スキームを受け取れるよう登録
      CFBundleURLTypes: reverseClientId ? [{ CFBundleURLSchemes: [reverseClientId] }] : [],
    },
  },
  android: {
    package: 'com.example.todoapp',
  },
  plugins: ['expo-router', 'expo-secure-store'],
}
