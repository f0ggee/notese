FROM golang:latest as builder

LABEL author = "fogge :0"

ENV DATABASE_URL="postgresql://postgres.swlbbzzazbagtprtcmgz:xokmix-tugweS-baqpa5@aws-0-eu-north-1.pooler.supabase.com:6543/postgres"
WORKDIR /app
EXPOSE 8080
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/main    main.go




FROM alpine:latest
WORKDIR /fogee
COPY --from=builder /app/main   /fogge/main
CMD ["fogge/main"]