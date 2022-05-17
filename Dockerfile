#
# swis-core-api / Dockerfile
#

# https://hub.docker.com/_/golang

ARG GOLANG_VERSION=1.17
FROM golang:${GOLANG_VERSION}-alpine

ARG APP_NAME
ARG DOCKER_DEV_PORT
ARG TZ

ENV APP_NAME ${APP_NAME}
ENV DOCKER_DEV_PORT ${DOCKER_DEV_PORT}
ENV TZ ${TZ}

RUN apk --no-cache add tzdata

WORKDIR /go/src/${APP_NAME}
COPY . .

RUN go mod init
RUN go get -d -v ./...
RUN go install 

EXPOSE ${DOCKER_DEV_PORT}
CMD ${APP_NAME}

