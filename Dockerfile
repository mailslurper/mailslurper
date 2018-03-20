FROM alpine:3.6

MAINTAINER erguotou525@gmail.compute

COPY . /tmp/mailslurper

ENV GOPATH /home/go
ENV GOBIN /home/go/bin
ENV PATH $PATH:$GOBIN
ENV ENABLE_CGO 1
ENV CGO_ENABLED 1

RUN \
  apk update \
  && apk add go git libc-dev \
  && mkdir -p /home/go/src/github.com/mailslurper /home/go/bin /home/go/pkg \
  && cd /home/go/src/github.com/mailslurper \
  && git clone https://github.com/mailslurper/mailslurper.git \
  && go get github.com/mjibson/esc \
  && cd mailslurper/cmd/mailslurper \
  && go get \
  && go generate \
  && go build \
  && rm /home/go/src/github.com/mailslurper/mailslurper/cmd/mailslurper/config.json

RUN echo -e '{\n\
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
  >> /home/go/src/github.com/mailslurper/mailslurper/cmd/mailslurper/config.json

WORKDIR /home/go/src/github.com/mailslurper/mailslurper/cmd/mailslurper

VOLUME /home/go/src/github.com/mailslurper/mailslurper

EXPOSE 8080 8085 2500

CMD /home/go/src/github.com/mailslurper/mailslurper/cmd/mailslurper/mailslurper
