FROM golang:alpine as builder

LABEL maintainer="erguotou525@gmail.compute"

RUN apk --no-cache add git libc-dev gcc

WORKDIR /app
COPY . /app
WORKDIR /go/src/github.com/mailslurper/mailslurper/cmd/mailslurper

RUN go get ./...
RUN go generate ./...
RUN go build ./cmd/mailslurper

FROM alpine:3.6

RUN apk add --no-cache ca-certificates \
 && echo -e '{\n\
  "wwwAddress": "0.0.0.0",\n\
  "wwwPort": 8080,\n\
  "serviceAddress": "0.0.0.0",\n\
  "servicePort": 8085,\n\
  "smtpAddress": "0.0.0.0",\n\
  "smtpPort": 2500,\n\
  "dbEngine": "SQLite",\n\
  "dbHost": "",\n\
  "dbPort": 0,\n\
  "dbDatabase": "./mailslurper.db",\n\
  "dbUserName": "",\n\
  "dbPassword": "",\n\
  "maxWorkers": 1000,\n\
  "autoStartBrowser": false,\n\
  "keyFile": "",\n\
  "certFile": "",\n\
  "adminKeyFile": "",\n\
  "adminCertFile": ""\n\
  }'\
  >> config.json

COPY --from=builder /app/mailslurper mailslurper

EXPOSE 8080 8085 2500

CMD ["./mailslurper"]
