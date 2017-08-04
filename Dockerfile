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
  && mkdir -p /home/go/src/github.com/mailslurper /home/go/bin /home/go/bin \
  && cd /home/go/src/github.com/mailslurper \
  && git clone https://github.com/mailslurper/libmailslurper.git \
  && mv /tmp/mailslurper ./ \
  && go get github.com/mjibson/esc \
  && cd mailslurper \
  && go get \
  && go generate \
  && go build

WORKDIR /home/go/src/github.com/mailslurper/mailslurper

VOLUME /home/go/src/github.com/mailslurper/mailslurper

EXPOSE 8080 8085 2500

CMD /home/go/src/github.com/mailslurper/mailslurper/mailslurper
