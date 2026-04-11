package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"

	"todo-app/backend/internal/auth"
	"todo-app/backend/internal/config"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func NewAuthInterceptor(cfg *config.Config) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// ExchangeToken と SignOut はスキップ（認証不要 / 期限切れでもログアウト可能）
			switch req.Spec().Procedure {
			case "/auth.v1.AuthService/ExchangeToken",
				"/auth.v1.AuthService/SignOut":
				return next(ctx, req)
			}

			// Authorization ヘッダ優先（Native）、なければ Cookie（Web）
			token := strings.TrimPrefix(req.Header().Get("Authorization"), "Bearer ")
			if token == "" {
				r := &http.Request{Header: req.Header()}
				if c, err := r.Cookie("jwt"); err == nil {
					token = c.Value
				}
			}

			userID, err := auth.ValidateJWT(cfg, token)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid or expired token"))
			}

			ctx = context.WithValue(ctx, UserIDKey, userID)
			return next(ctx, req)
		}
	}
}
