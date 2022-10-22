# swis-api (swapi)

[![swis-api CI/CD pipeline](https://github.com/savla-dev/swis-api/actions/workflows/docker-image.yml/badge.svg)](https://github.com/savla-dev/swis-api/actions/workflows/docker-image.yml)

[sakalWeb (v5) Information System RESTful API (intranet)](http://swapi.savla.su)

## repo vademecum

### .assets

This folder contains static assets like images (favicon), that are served by server statically (not dynamically).

### .script

This folder contains executables for swapi package management, like data importing batch script and backuping script(s).

### .env / .env.example

Important configuration file containing environment constants for `Makefile` and other linked files like `docker-compose.yml`. Some constants, like `ROOT_TOKEN` are essential for the app startup.

### main.go

Main package's logic is coded there, main package's entrypoint, server router and its route groups (entrypoints for other modules) are present there as well.

### modules

All swapi modules are stored in their folders. Every module has its `models.go` file with data structures, `controllers.go` with its methods and functions, and `routes.go` for gin router to serve module's handles.

More on swis-api modules in [this article](https://krusty.savla.dev/projects/swis-api/).

### .docker/Dockerfile

Recipe for docker image build using `docker build .`

### Makefile

Dev/build recipe for GNU `make` tool.

```shell
$ make

 swis-api / Makefile 

 make --- show this helper 

 make fmt  --- reformat the go source (gofmt) 
 make docs --- render documentation from code (swagger OA docs) 

 make build --- build project (docker image) 
 make run   --- run project 
 make logs  --- fetch container's logs 
 make stop  --- stop and purge project (only docker containers!) 

 make import_dump --- import dumped data (locally) into runtime 
 make dump        --- dump runtime data to DUMP_DIR 
 make backup      --- execute data dump and tar/gzip data backup 

```

### docker-compose.yml

YAML-formated file for docker-compose stack. Contains defitions for docker container and its isolated network. Uses constants from `.env` dot-file.


## documentation

[swagger 2.0 is used to document API scheme (priv)](http://swapi-docs.savla.su)

```
# generate docs using swaggo/swag and restart swagger_ui container
make docs

# view swagger UI
http://localhost:8999/
```

## staging and deployment

As far as the deployment architecture is concerned, the codebase is designed to follow (mostly) [12factor guidelines](https://12factor.net). Thus there is a difference between a build and release and runtime stages.

![swis-api-pipeline](./.assets/swis-api-pipeline.png)

## authentication

Swapi uses token-based authentication for any request to be authenticated, For initial importing, `ROOT_TOKEN` (see [`.env`](/.env) file) is used by importing executable. For any request, the header `X-Auth-Token` has to be sent with a custom HTTP request.

CORS is set only for http://swjango.savla.su at the moment.

## service backup report example

```shell
SIZE=$(du -shx ${BACKUP_TARGET_DIR}/${TIMESTAMP}.sqlite.gz | awk '{ print $1 }')
STATUS=success

# report back to swapi/backups
SERVICE_NAME=generic-sqlite-service
TIMESTAMP=$(date +%s)
TOKEN=xxx
curl -X PUT -sL -H "X-Auth-Token: $TOKEN" \
        --data "{\"service_name\":\"$SERVICE_NAME\", \"timestamp\": $TIMESTAMP, \"last_status\": \"$STATUS\", \"backup_size\": \"$SIZE\", \"filename\": \"${TIMESTAMP}.sqlite.gz\" }" \
        http://swis-api-hostname/backups/${SERVICE_NAME}
```

## importing

At start, swapi instance memory is cleared and ready for any data import (until the next restart). Any data stored in runtime memory should be dumped using GET methods at particular paths. This approach should make `swapi` instance universal (while omitting custom packages/modules).

```shell
# build new image version
make build

# dump production data locally
make dump

# redeploy the running server with new image
make run

# import prod data from local DUMP_DIR (see .env)
make import_dump
```

## roadmap for v5.3 (public release)

+ rebase master branch (yeet sensitive commits)
+ deep refactor -- omit global vars, clean the code
+ ~~use muxers -- save/dump runtime in-memory data~~
+ improve backuping -- introduce tiny go binary to backup all the data
+ add CRUD controllers to each module/package

### nice-to-have(s)

+ try to implement generics/first class functions
+ implement 2FA (but only do store Google Auth keys) for frontend
+ improve data structure by using maps
+ auto-data-recovery after reboot
