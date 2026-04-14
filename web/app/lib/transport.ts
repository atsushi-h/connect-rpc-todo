import { createConnectTransport } from "@connectrpc/connect-web";

export const transport = createConnectTransport({
  baseUrl: import.meta.env.VITE_API_URL ?? "http://localhost:8080",
  // HttpOnly Cookie を自動送信するために credentials: "include" を fetch でオーバーライド
  fetch: (input, init) => fetch(input, { ...init, credentials: "include" }),
});
