# Стадия сборки
FROM golang:latest AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY fronted ./fronted
COPY cmd ./cmd
COPY iteranal ./iteranal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./





FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/app .
COPY --from=builder /build/fronted ./fronted
COPY --from=builder /build/cmd ./cmd

COPY .env .env
EXPOSE 443



CMD ["./app"]
