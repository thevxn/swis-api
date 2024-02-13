#
# swis-api / Dockerfile
#

#
# stage 0 -- build image executable, combine codes and resource links
#

# https://hub.docker.com/_/golang
ARG GOLANG_VERSION
ARG ALPINE_VERSION
FROM golang:${GOLANG_VERSION}-alpine AS swapi-build

ARG APP_NAME
ENV APP_NAME ${APP_NAME}

RUN apk add git

WORKDIR /go/src/${APP_NAME}
COPY ./go.mod .
RUN go mod download

COPY . .

# force use locked version
ARG APP_VERSION
#RUN go install go.savla.dev/swis/v5@v${APP_VERSION}
RUN go build -o /go/bin/swis


#
# stage 1 -- release, combine build image and config(s), pre-runtime
#

FROM alpine:${ALPINE_VERSION} AS swapi-release

ARG APP_NAME
ARG APP_ROOT
ARG DOCKER_INTERNAL_PORT
ARG DOCKER_USER

ENV APP_NAME ${APP_NAME}

RUN apk add tzdata
RUN adduser -D -h ${APP_ROOT} -s /bin/sh ${DOCKER_USER}

# compose release image
COPY --from=swapi-build /go/bin/swis /bin/${APP_NAME}
COPY ./.github/assets/favicon.ico ${APP_ROOT}/

WORKDIR ${APP_ROOT}
EXPOSE ${DOCKER_INTERNAL_PORT}
USER ${DOCKER_USER}
ENTRYPOINT ["/bin/swis-api"]

