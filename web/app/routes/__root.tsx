import { createRootRouteWithContext, Outlet, Scripts, HeadContent } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { TransportProvider } from "@connectrpc/connect-query";
import { createClient } from "@connectrpc/connect";
import { createGrpcWebTransport } from "@connectrpc/connect-web";
import { AuthService } from "@todo-app/api-client/src/auth/v1/auth_pb.js";
import { transport } from "~/lib/transport";
import type { ReactNode } from "react";

interface RouterContext {
  queryClient: QueryClient;
}

export interface AuthContext {
  isAuthenticated: boolean;
  user?: {
    id: string;
    email: string;
    displayName: string;
    avatarUrl: string;
  };
}

const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

function NotFound() {
  return <p>Not Found</p>;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  notFoundComponent: NotFound,
  beforeLoad: async (): Promise<{ auth: AuthContext }> => {
    let authTransport = transport;

    if (typeof window === "undefined") {
      // SSR: 受信リクエストの Cookie をバックエンドへ転送する
      const { getRequestHeader } = await import("@tanstack/react-start/server");
      const cookieHeader = getRequestHeader("cookie") ?? "";
      authTransport = createGrpcWebTransport({
        baseUrl: API_BASE,
        fetch: (input, init) => {
          const headers = new Headers(init?.headers);
          if (cookieHeader) headers.set("cookie", cookieHeader);
          return fetch(input as string, { ...(init as RequestInit), headers });
        },
      });
    }

    const client = createClient(AuthService, authTransport);
    try {
      const me = await client.getMe({});
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
      };
    } catch {
      return { auth: { isAuthenticated: false } };
    }
  },
  component: RootComponent,
});

function RootDocument({ children }: { children: ReactNode }) {
  return (
    <html>
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <Scripts />
      </body>
    </html>
  );
}

function RootComponent() {
  const { queryClient } = Route.useRouteContext();
  return (
    <RootDocument>
      <TransportProvider transport={transport}>
        <QueryClientProvider client={queryClient}>
          <Outlet />
        </QueryClientProvider>
      </TransportProvider>
    </RootDocument>
  );
}
