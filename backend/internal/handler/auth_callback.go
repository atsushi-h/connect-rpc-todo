package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"todo-app/backend/internal/auth"
	"todo-app/backend/internal/config"
	"todo-app/backend/internal/db"
)

type AuthCallbackHandler struct {
	queries     *db.Queries
	cfg         *config.Config
	oauthConfig *oauth2.Config
}

func NewAuthCallbackHandler(queries *db.Queries, cfg *config.Config) *AuthCallbackHandler {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleCallbackURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	return &AuthCallbackHandler{queries: queries, cfg: cfg, oauthConfig: oauthConfig}
}

func (h *AuthCallbackHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := generateState()
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   600,
		})

		url := h.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func (h *AuthCallbackHandler) Callback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// state 検証
		stateCookie, err := r.Cookie("oauth_state")
		if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:   "oauth_state",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		// authorization code を access token に交換
		token, err := h.oauthConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			log.Printf("ERROR callback exchange: %v", err)
			http.Error(w, "failed to exchange token", http.StatusInternalServerError)
			return
		}

		// Google userinfo 取得
		userInfo, err := fetchGoogleUserInfo(r.Context(), h.oauthConfig, token)
		if err != nil {
			log.Printf("ERROR callback userinfo: %v", err)
			http.Error(w, "failed to get user info", http.StatusInternalServerError)
			return
		}

		// users テーブルに upsert
		user, err := h.queries.UpsertUserByGoogleID(r.Context(), db.UpsertUserByGoogleIDParams{
			Email:       userInfo.Email,
			GoogleID:    sql.NullString{String: userInfo.Sub, Valid: true},
			DisplayName: userInfo.Name,
			AvatarUrl:   userInfo.Picture,
		})
		if err != nil {
			log.Printf("ERROR callback upsert: %v", err)
			http.Error(w, "failed to upsert user", http.StatusInternalServerError)
			return
		}

		// JWT 生成
		jwt, err := auth.SignJWT(h.cfg, user.ID.String())
		if err != nil {
			log.Printf("ERROR callback sign jwt: %v", err)
			http.Error(w, "failed to create token", http.StatusInternalServerError)
			return
		}

		// JWT を HttpOnly Cookie にセット
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    jwt,
			Path:     "/",
			HttpOnly: true,
			Secure:   h.cfg.CookieSecure,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   30 * 24 * 60 * 60,
		})

		http.Redirect(w, r, h.cfg.WebFrontendURL, http.StatusFound)
	}
}

type googleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func fetchGoogleUserInfo(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (*googleUserInfo, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google userinfo: unexpected status %d", resp.StatusCode)
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	if info.Sub == "" || info.Email == "" {
		return nil, fmt.Errorf("google userinfo: missing required fields")
	}
	return &info, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
