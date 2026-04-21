package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"todo-app/backend/internal/auth"
	"todo-app/backend/internal/config"
)

func testConfig(secret string) *config.Config {
	return &config.Config{JWTSecret: secret}
}

func TestSignAndValidateJWT(t *testing.T) {
	t.Parallel()

	cfg := testConfig("test-secret")
	userID := "user-123"

	token, err := auth.SignJWT(cfg, userID)
	if err != nil {
		t.Fatalf("SignJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("SignJWT returned empty token")
	}

	got, err := auth.ValidateJWT(cfg, token)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if got != userID {
		t.Errorf("ValidateJWT userID = %q, want %q", got, userID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	t.Parallel()

	cfg := testConfig("test-secret")

	claims := auth.Claims{
		UserID: "user-expired",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = auth.ValidateJWT(cfg, tokenStr)
	if err == nil {
		t.Fatal("ValidateJWT should fail for expired token, but got nil error")
	}
}

func TestValidateJWT_TamperedToken(t *testing.T) {
	t.Parallel()

	cfg := testConfig("test-secret")

	token, err := auth.SignJWT(cfg, "user-123")
	if err != nil {
		t.Fatalf("SignJWT failed: %v", err)
	}

	// 末尾1文字を変えて改ざん
	tampered := token[:len(token)-1] + "X"

	_, err = auth.ValidateJWT(cfg, tampered)
	if err == nil {
		t.Fatal("ValidateJWT should fail for tampered token, but got nil error")
	}
}

func TestValidateJWT_WrongSigningMethod(t *testing.T) {
	t.Parallel()

	cfg := testConfig("test-secret")

	// テスト専用 RSA 鍵を生成して RS256 で署名
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	claims := auth.Claims{
		UserID: "user-rsa",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := token.SignedString(rsaKey)
	if err != nil {
		t.Fatalf("failed to sign with RSA: %v", err)
	}

	_, err = auth.ValidateJWT(cfg, tokenStr)
	if err == nil {
		t.Fatal("ValidateJWT should fail for RS256-signed token, but got nil error")
	}
}

func TestValidateJWT_EmptyToken(t *testing.T) {
	t.Parallel()

	cfg := testConfig("test-secret")

	_, err := auth.ValidateJWT(cfg, "")
	if err == nil {
		t.Fatal("ValidateJWT should fail for empty token, but got nil error")
	}
}
