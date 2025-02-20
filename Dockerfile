FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY ["tg/", "tg/"]
COPY ["main.go/", "main.go"]
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/app main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/app ./app
COPY ["config.json", "config.json"]

ENTRYPOINT [ "/app" ]
