FROM golang:1.26.2 AS builder

WORKDIR /app

# go.work で backend と gen/go を束ねているので両方コピー
COPY go.work go.work
COPY backend/go.mod backend/go.sum ./backend/
COPY gen/go/ ./gen/go/

RUN cd backend && go mod download

COPY backend/ ./backend/

RUN cd backend && go build -o /server ./cmd/server

FROM gcr.io/distroless/base-debian12
COPY --from=builder /server /server
EXPOSE 8080
CMD ["/server"]
