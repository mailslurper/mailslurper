FROM golang:alpine as builder



LABEL maintainer="erguotou525@gmail.compute"

RUN apk --no-cache add git libc-dev gcc
RUN go install github.com/mjibson/esc@latest

COPY . /go/src/github.com/mailslurper/mailslurper
WORKDIR /go/src/github.com/mailslurper/mailslurper/cmd/mailslurper

RUN go install
RUN go generate
RUN go build

FROM alpine:3.6
ARG wwwPort=8080
ARG wwwAddress=0.0.0.0
ARG wwwPublicURL
ARG serviceAddress=0.0.0.0
ARG servicePort=8005
ARG servicePublicURL
ARG smtpAddress=0.0.0.0
ARG smtpPort=2500
ARG dbEngine=SQLite
ARG dbHost
ARG dbPort=0
ARG dbDatabase=".\/mailslurper"
ARG dbUserName
ARG dbPassword

RUN apk add --no-cache ca-certificates \
 && echo -e '{\n\
  "wwwAddress": "www_Address",\n\
  "wwwPort": www_Port,\n\
  "wwwPublicURL": "www_PublicURL",\n\
  "serviceAddress": "service_Address",\n\
  "servicePort": service_Port,\n\
  "servicePublicURL": "service_PublicURL",\n\
  "smtpAddress": "smtp_Address",\n\
  "smtpPort": smtp_Port,\n\
  "dbEngine": "db_Engine",\n\
  "dbHost": "db_Host",\n\
  "dbPort": db_Port,\n\
  "dbDatabase": "db_Database",\n\
  "dbUserName": "db_UserName",\n\
  "dbPassword": "db_Password",\n\
  "maxWorkers": 1000,\n\
  "autoStartBrowser": false,\n\
  "keyFile": "",\n\
  "certFile": "",\n\
  "adminKeyFile": "",\n\
  "adminCertFile": ""\n\
  }'\
  >> config.json 
 
 RUN sed "s/www_Port/${wwwPort}/g" config.json > config1.json && \
 sed "s/www_Address/${wwwAddress}/g" config1.json > config2.json && \
 sed "s/www_PublicURL/${wwwPublicURL}/g" config2.json > config3.json && \
 sed "s/service_Address/${serviceAddress}/g" config3.json > config4.json && \
 sed "s/service_Port/${servicePort}/g" config4.json > config5.json && \
 sed "s/service_PublicURL/${servicePublicURL}/g" config5.json > config6.json && \
 sed "s/smtp_Address/${smtpAddress}/g" config6.json > config7.json && \
 sed "s/smtp_Port/${smtpPort}/g" config7.json > config8.json && \
 sed "s/db_Engine/${dbEngine}/g" config8.json > config9.json && \
 sed "s/db_Host/${dbHost}/g" config9.json > config10.json && \
 sed "s/db_Port/${dbPort}/g" config10.json > config11.json && \
 sed "s/db_Database/${dbDatabase}/g" config11.json > config12.json && \
 sed "s/db_UserName/${dbUserName}/g" config12.json > config13.json && \
 sed "s/db_Password/${dbPassword}/g" config13.json > config.json && \
 rm -rf config1.json && \
 rm -rf config2.json && \
 rm -rf config3.json && \
 rm -rf config4.json && \
 rm -rf config5.json && \
 rm -rf config6.json && \
 rm -rf config7.json && \
 rm -rf config8.json && \
 rm -rf config9.json && \
 rm -rf config10.json && \
 rm -rf config11.json && \
 rm -rf config12.json && \
 rm -rf config13.json

COPY --from=builder /go/src/github.com/mailslurper/mailslurper/cmd/mailslurper/mailslurper mailslurper

EXPOSE ${wwwPort} ${servicePort} ${smtpPort}

CMD ["./mailslurper"]
