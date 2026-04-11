package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	authv1 "gen/go/auth/v1"
	"gen/go/auth/v1/authv1connect"
	"todo-app/backend/internal/auth"
	"todo-app/backend/internal/config"
	"todo-app/backend/internal/db"
	"todo-app/backend/internal/middleware"
)

type AuthHandler struct {
	authv1connect.UnimplementedAuthServiceHandler
	queries     *db.Queries
	cfg         *config.Config
	oauthConfig *oauth2.Config
}

func NewAuthHandler(queries *db.Queries, cfg *config.Config) *AuthHandler {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	return &AuthHandler{queries: queries, cfg: cfg, oauthConfig: oauthConfig}
}

// ExchangeToken は Native PKCE フロー用: authorization code を JWT に交換する
func (h *AuthHandler) ExchangeToken(
	ctx context.Context,
	req *connect.Request[authv1.ExchangeTokenRequest],
) (*connect.Response[authv1.ExchangeTokenResponse], error) {
	cfg := *h.oauthConfig
	cfg.RedirectURL = req.Msg.RedirectUri

	token, err := cfg.Exchange(ctx, req.Msg.Code,
		oauth2.SetAuthURLParam("code_verifier", req.Msg.CodeVerifier),
	)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid authorization code"))
	}

	userInfo, err := fetchGoogleUserInfo(ctx, &cfg, token)
	if err != nil {
		log.Printf("ERROR ExchangeToken userinfo: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
	}

	user, err := h.queries.UpsertUserByGoogleID(ctx, db.UpsertUserByGoogleIDParams{
		Email:       userInfo.Email,
		GoogleID:    sql.NullString{String: userInfo.Sub, Valid: true},
		DisplayName: userInfo.Name,
		AvatarUrl:   userInfo.Picture,
	})
	if err != nil {
		log.Printf("ERROR ExchangeToken upsert: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
	}

	jwt, err := auth.SignJWT(h.cfg, user.ID.String())
	if err != nil {
		log.Printf("ERROR ExchangeToken sign jwt: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
	}

	return connect.NewResponse(&authv1.ExchangeTokenResponse{AccessToken: jwt}), nil
}

// GetMe は認証済みユーザーの情報を返す（Web / Native 共通）
func (h *AuthHandler) GetMe(
	ctx context.Context,
	req *connect.Request[authv1.GetMeRequest],
) (*connect.Response[authv1.GetMeResponse], error) {
	userIDStr, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("ERROR GetMe parse uuid: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
	}

	user, err := h.queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
		}
		log.Printf("ERROR GetMe db: %v", err)
		return nil, connect.NewError(connect.CodeInternal, errors.New("internal server error"))
	}

	return connect.NewResponse(&authv1.GetMeResponse{
		Id:          user.ID.String(),
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarUrl:   user.AvatarUrl,
	}), nil
}

// SignOut は Web 向けに JWT Cookie を削除する
func (h *AuthHandler) SignOut(
	ctx context.Context,
	req *connect.Request[authv1.SignOutRequest],
) (*connect.Response[authv1.SignOutResponse], error) {
	resp := connect.NewResponse(&authv1.SignOutResponse{})
	resp.Header().Add("Set-Cookie", (&http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}).String())
	return resp, nil
}
