FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY . .

# No CGO needed: modernc.org/sqlite is pure Go, and the frontend is plain
# static files embedded via go:embed, so no separate asset-bundling step
# (esc, go generate) is required either.
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/mylslurper ./cmd/mylslurper
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/createcredentials ./cmd/createcredentials

FROM alpine:latest AS certs

RUN apk add --no-cache openssl \
 && mkdir -p /certs \
 && openssl req -x509 -newkey rsa:2048 -nodes \
      -keyout /certs/server.key \
      -out /certs/server.crt \
      -days 3650 \
      -subj "/CN=localhost" \
      -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

FROM alpine:latest

RUN apk add --no-cache ca-certificates wget

COPY --from=builder /out/mylslurper /out/createcredentials /app/
COPY --from=certs /certs/server.crt /certs/server.key /app/
WORKDIR /app

# Implicit TLS on SMTP (smtps://) and HTTPS on the web UI. Unset both
# to fall back to plaintext.
ENV CERT_FILE=/app/server.crt \
    KEY_FILE=/app/server.key

EXPOSE 4436 4437 1025

ENTRYPOINT ["/app/mylslurper"]
