-- name: UpsertUserByGoogleID :one
INSERT INTO users (email, google_id, display_name, avatar_url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (google_id) DO UPDATE
  SET email        = EXCLUDED.email,
      display_name = EXCLUDED.display_name,
      avatar_url   = EXCLUDED.avatar_url,
      updated_at   = now()
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;
