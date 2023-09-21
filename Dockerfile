FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.19.13-alpine3.18 as builder

ARG TARGETOS TARGETARCH

LABEL maintainer="erguotou525@gmail.compute"

RUN apk --no-cache add git libc-dev gcc ca-certificates
RUN go install github.com/mjibson/esc@v0.2.0

RUN ls "$GOPATH/bin"

COPY . /go/src/github.com/mailslurper/mailslurper
WORKDIR /go/src/github.com/mailslurper/mailslurper/cmd/mailslurper

ENV GOOS=$TARGETOS GOARCH=$TARGETARCH

RUN go get
RUN go generate
RUN go build

RUN <<EOF
echo -e '{
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
}' >> config.json
EOF

FROM gcr.io/distroless/static-debian12

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /go/src/github.com/mailslurper/mailslurper/cmd/mailslurper/mailslurper mailslurper

EXPOSE 8080 8085 2500

CMD ["./mailslurper"]
