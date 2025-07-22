FROM alpine:latest AS ffmpeg
# Install curl to download the static FFmpeg binary
RUN apk add --no-cache curl

# Download and extract the static FFmpeg build
RUN curl -LO https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz \
    && tar xf ffmpeg-release-amd64-static.tar.xz \
    && mv ffmpeg-*-static/ffmpeg /ffmpeg \
    && mv ffmpeg-*-static/ffprobe /ffprobe

FROM golang:1.23-alpine AS go
WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY ["tg/", "tg/"]
COPY ["main.go/", "main.go"]
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/app main.go

FROM scratch

COPY --from=ffmpeg /tmp /tmp
COPY --from=ffmpeg /ffmpeg /usr/local/bin/ffmpeg
COPY --from=ffmpeg /ffprobe /usr/local/bin/ffprobe
COPY --from=ffmpeg /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go /app/bin/app ./app
COPY ["config.json", "config.json"]

ENTRYPOINT [ "/app" ]
