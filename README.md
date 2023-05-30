# swis-api (swapi) v5.4

[![Go Reference](https://pkg.go.dev/badge/go.savla.dev/swis/v5.svg)](https://pkg.go.dev/go.savla.dev/swis/v5)
[![Go Report Card](https://goreportcard.com/badge/go.savla.dev/swis/v5)](https://goreportcard.com/report/go.savla.dev/swis/v5)

sakalWebIS v5 RESTful JSON API in Go for Docker

`swapi` is a condensated data structure tree made into an Information System (IS) environment. It works as an in-memory database/cache for various system components like `users`, `roles`, `projects`, `business_entities` and more. Primarily, it is meant to be used in Docker as all components are written to work out of the box when there.

Reaad more in a short article on `swapi`:

+ https://krusty.space/projects/swis-api

### simple way to run (docker)

```
git clone https://github.com/savla-dev/swis-api
cp .env.example .env
vi .env

make build run
```

### install binary

```
go install go.savla.dev/swis/v5@latest
```

`swapi` could be run as a single binary too. However, some environment constants have to be set when running it solitarily --- mainly `DOCKER_INTERNAL_PORT` and `ROOT_TOKEN` constants are required for the app smooth start. 
