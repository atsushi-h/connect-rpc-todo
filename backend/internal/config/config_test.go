package config_test

import (
	"os"
	"testing"

	"todo-app/backend/internal/config"
)

// clearEnv は config.Load() が参照する全 env vars をクリアして
// テスト終了後に元の値を復元する。
func clearEnv(t *testing.T) {
	t.Helper()
	keys := []string{
		"DATABASE_URL", "SERVER_PORT", "JWT_SECRET",
		"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "GOOGLE_NATIVE_CLIENT_ID",
		"GOOGLE_CALLBACK_URL", "WEB_FRONTEND_URL", "COOKIE_SECURE",
	}
	saved := make(map[string]string, len(keys))
	for _, k := range keys {
		saved[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for k, v := range saved {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})
}

// setRequired は Load() に必要な最低限の env vars を設定する。
func setRequired(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost/testdb")
	t.Setenv("JWT_SECRET", "test-jwt-secret")
	t.Setenv("GOOGLE_CLIENT_ID", "test-google-client-id")
}

func TestLoad_Success(t *testing.T) {
	clearEnv(t)
	setRequired(t)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.DatabaseURL != "postgres://user:pass@localhost/testdb" {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, "postgres://user:pass@localhost/testdb")
	}
	if cfg.JWTSecret != "test-jwt-secret" {
		t.Errorf("JWTSecret = %q, want %q", cfg.JWTSecret, "test-jwt-secret")
	}
	if cfg.GoogleClientID != "test-google-client-id" {
		t.Errorf("GoogleClientID = %q, want %q", cfg.GoogleClientID, "test-google-client-id")
	}
}

func TestLoad_DefaultServerPort(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	// SERVER_PORT は未設定のまま

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.ServerPort != 8080 {
		t.Errorf("ServerPort = %d, want 8080 (default)", cfg.ServerPort)
	}
}

func TestLoad_CustomServerPort(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	t.Setenv("SERVER_PORT", "9090")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.ServerPort != 9090 {
		t.Errorf("ServerPort = %d, want 9090", cfg.ServerPort)
	}
}

func TestLoad_InvalidServerPort(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	t.Setenv("SERVER_PORT", "not-a-number")

	_, err := config.Load()
	if err == nil {
		t.Fatal("Load should fail for invalid SERVER_PORT, but got nil error")
	}
}

func TestLoad_CookieSecureTrue(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	t.Setenv("COOKIE_SECURE", "true")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if !cfg.CookieSecure {
		t.Error("CookieSecure should be true when COOKIE_SECURE=true")
	}
}

func TestLoad_CookieSecureFalse(t *testing.T) {
	clearEnv(t)
	setRequired(t)
	// COOKIE_SECURE は未設定のまま

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.CookieSecure {
		t.Error("CookieSecure should be false when COOKIE_SECURE is not set")
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	clearEnv(t)
	t.Setenv("JWT_SECRET", "test-jwt-secret")
	t.Setenv("GOOGLE_CLIENT_ID", "test-google-client-id")
	// DATABASE_URL は未設定のまま

	_, err := config.Load()
	if err == nil {
		t.Fatal("Load should fail when DATABASE_URL is not set, but got nil error")
	}
}

func TestLoad_MissingJWTSecret(t *testing.T) {
	clearEnv(t)
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost/testdb")
	t.Setenv("GOOGLE_CLIENT_ID", "test-google-client-id")
	// JWT_SECRET は未設定のまま

	_, err := config.Load()
	if err == nil {
		t.Fatal("Load should fail when JWT_SECRET is not set, but got nil error")
	}
}

func TestLoad_MissingGoogleClientID(t *testing.T) {
	clearEnv(t)
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost/testdb")
	t.Setenv("JWT_SECRET", "test-jwt-secret")
	// GOOGLE_CLIENT_ID は未設定のまま

	_, err := config.Load()
	if err == nil {
		t.Fatal("Load should fail when GOOGLE_CLIENT_ID is not set, but got nil error")
	}
}
