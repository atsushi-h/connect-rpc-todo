# CLAUDE.md

## プロジェクト概要

Connect RPC を使った認証機能付き Todo リストアプリ。Web / Native でフロントエンドを共有するモノレポ構成。

## 技術スタック

- **スキーマ定義:** Protocol Buffers + Buf
- **Backend:** Go, connect-go, sqlc, golang-migrate
- **Frontend (Web):** TanStack Start, connect-query
- **Frontend (Native):** React Native (Expo), connect-query
- **DB:** PostgreSQL 16
- **ランタイム管理:** mise (Go, Node, pnpm)
- **パッケージ管理:** pnpm workspace, Go workspace
- **コンテナ:** Docker Compose

## ディレクトリ構成

```
todo-app/
├── .mise.toml                 # Go, Node, pnpm バージョン固定
├── go.work                    # Go workspace（backend + gen/go）
├── package.json               # pnpm workspace root
├── pnpm-workspace.yaml
├── compose.yaml               # PostgreSQL
├── proto/                     # .proto 定義（単一の信頼源）
│   ├── buf.yaml
│   ├── buf.gen.yaml
│   └── todo/v1/
│       └── todo.proto
├── gen/                       # buf generate の出力先
│   ├── go/                    #   → backend が参照
│   └── ts/                    #   → web / native が workspace パッケージとして参照
├── backend/                   # Go module
│   ├── go.mod
│   ├── cmd/server/main.go
│   └── internal/
│       ├── handler/           # Connect RPC ハンドラ
│       ├── middleware/         # interceptor（認証等）
│       ├── repository/
│       └── db/                # sqlc 生成コード
├── web/                       # TanStack Start
│   └── package.json
└── native/                    # Expo (React Native)
    ├── package.json
    └── metro.config.js
```

## よく使うコマンド

```bash
# セットアップ
mise install
pnpm install
docker compose up -d db

# Proto からコード生成
cd proto && buf lint && buf generate

# Backend 起動
cd backend && go run ./cmd/server

# Web 起動
cd web && pnpm dev

# Native 起動
cd native && pnpm expo start

# DB マイグレーション
cd backend && migrate -path db/migrations -database "postgres://todo:todo@localhost:5432/todo?sslmode=disable" up

# sqlc 再生成
cd backend && sqlc generate
```