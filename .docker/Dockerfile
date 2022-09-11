#
# swis-core-api / Dockerfile
#

#
# stage 0 -- build image executable, combine codes and resource links
#

# https://hub.docker.com/_/golang
ARG GOLANG_VERSION
ARG ALPINE_VERSION
FROM golang:${GOLANG_VERSION}-alpine as swapi-build

ARG APP_NAME
ENV APP_NAME ${APP_NAME}

RUN apk --no-cache add git

WORKDIR /go/src/${APP_NAME}
COPY . .

RUN go mod init ${APP_NAME}
RUN go mod tidy

# swagger documentation
RUN go get -u github.com/swaggo/swag/cmd/swag && \
	go install github.com/swaggo/swag/cmd/swag@latest && \
	swag init .

RUN go install ${APP_NAME}


#
# stage 1 -- release, combine build image and config(s)
#

FROM alpine:${ALPINE_VERSION} as swapi-release

ARG APP_FLAGS
ARG APP_NAME
ARG APP_ROOT
ARG DOCKER_DEV_PORT
ARG DOCKER_USER
ARG GIN_MODE
ARG ROOT_TOKEN
ARG TZ

ENV APP_FLAGS ${APP_FLAGS}
ENV APP_NAME ${APP_NAME}
ENV APP_ROOT ${APP_ROOT}
ENV GIN_MODE ${GIN_MODE}
ENV ROOT_TOKEN ${ROOT_TOKEN}
ENV TZ ${TZ}

RUN apk --no-cache add tzdata
RUN adduser -D -h ${APP_ROOT} -s /bin/sh ${DOCKER_USER}

# compose release image
COPY --from=swapi-build /go/bin/${APP_NAME} /bin/${APP_NAME}

WORKDIR ${APP_ROOT}
EXPOSE ${DOCKER_APP_PORT}
USER ${DOCKER_USER}
ENTRYPOINT ${APP_NAME}
CMD ${APP_FLAGS}

