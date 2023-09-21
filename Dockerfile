FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.21.1-bookworm as builder

ARG ESC_VERSION="v0.2.0"

LABEL maintainer="erguotou525@gmail.compute"

ARG TARGETOS TARGETARCH

RUN dpkg --add-architecture "${TARGETARCH}"
RUN apt update
RUN <<EOF
case "${TARGETARCH}" in
  "amd64")
    CC_PACKAGE="gcc-x86-64-linux-gnu" ;;
  "arm64")
    CC_PACKAGE="gcc-aarch64-linux-gnu" ;;
esac
apt install --yes "$CC_PACKAGE" git libc-dev ca-certificates libsqlite3-dev:"${TARGETARCH}"
EOF

RUN go install github.com/mjibson/esc@"${ESC_VERSION}"

COPY . /go/src/github.com/mailslurper/mailslurper
WORKDIR /go/src/github.com/mailslurper/mailslurper/cmd/mailslurper

ENV GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" CGO_ENABLED=1

RUN go get
RUN go generate
RUN <<EOF
case "${TARGETARCH}" in
  "amd64")
    CC_NAME="x86_64-linux-gnu-gcc" ;;
  "arm64")
    CC_NAME="aarch64-linux-gnu-gcc" ;;
esac
CC="$CC_NAME" go build -o /out/mailslurper
EOF
RUN

RUN <<EOF
echo '{
  "wwwAddress": "0.0.0.0",
  "wwwPort": 8080,
  "wwwPublicURL": "",
  "serviceAddress": "0.0.0.0",
  "servicePort": 8085,
  "servicePublicURL": "",
  "smtpAddress": "0.0.0.0",
  "smtpPort": 2500,
  "dbEngine": "SQLite",
  "dbHost": "",
  "dbPort": 0,
  "dbDatabase": "./mailslurper.db",
  "dbUserName": "",
  "dbPassword": "",
  "maxWorkers": 1000,
  "autoStartBrowser": false,
  "keyFile": "",
  "certFile": "",
  "adminKeyFile": "",
  "adminCertFile": ""
}' > /out/config.json
EOF

FROM gcr.io/distroless/base-debian12

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /out/mailslurper /bin/mailslurper
COPY --from=builder /out/config.json /etc/mailslurper/config.json

EXPOSE 8080 8085 2500

CMD ["/bin/mailslurper", "--config", "/etc/mailslurper/config.json"]
