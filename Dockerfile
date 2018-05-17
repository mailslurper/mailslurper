FROM golang:alpine as builder

LABEL maintainer="erguotou525@gmail.compute"

RUN apk --no-cache add git libc-dev gcc
RUN go get github.com/mjibson/esc

COPY . ./src

RUN cd src/cmd/mailslurper \
 && go get \
 && go generate \
 && go build


FROM alpine:3.6


COPY --from=builder /go/src/cmd/mailslurper/mailslurper mailslurper

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

EXPOSE 8080 8085 2500

CMD ["./mailslurper"]
