# Стадия сборки
FROM golang:latest AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY fronted/ .
COPY iteranal/ .
COPY cmd /.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /main  /app/main
COPY --from=builder /build/fronted  /app/fronted

COPY .env  /app/.env
ENV DATABASE_URL=postgresql://postgres.swlbbzzazbagtprtcmgz:xokmix-tugweS-baqpa5@aws-0-eu-north-1.pooler.supabase.com:6543/postgres

EXPOSE 8080
CMD ["/app/main"]