FROM golang:1.11-alpine3.9 AS build

MAINTAINER Samuele Catuzzi

ENV REPO github.com/samuelecat/oauthbridge

RUN apk add git bash

RUN go get golang.org/x/oauth2
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/yaml.v2

#RUN go get ${REPO}
# uncomment the following if you are working locally
WORKDIR /go/src/${REPO}
COPY . .
RUN go install

FROM alpine:3.9
RUN apk add --update --no-cache \
		ca-certificates \
		openssl \
		nginx \
		supervisor

COPY --from=build /go/bin/oauthbridge /usr/bin/
COPY Dockerfiles/supervisord.conf /etc/
COPY Dockerfiles/nginx.ini /etc/supervisor.d/
COPY Dockerfiles/oauthbridge.ini /etc/supervisor.d/

COPY Dockerfiles/nginx-server.conf /etc/nginx/conf.d/
RUN mkdir /run/nginx

VOLUME /etc/oauthbridge 

# 8081 nginx reverse proxy, 9999 oauthbridge, 9001 supervisord
EXPOSE 8081 9999 9001

CMD ["/usr/bin/supervisord", "-n"]
