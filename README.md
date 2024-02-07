# swis-api (swapi) v5.12

[![Go Reference](https://pkg.go.dev/badge/go.savla.dev/swis/v5.svg)](https://pkg.go.dev/go.savla.dev/swis/v5)
[![Go Report Card](https://goreportcard.com/badge/go.savla.dev/swis/v5)](https://goreportcard.com/report/go.savla.dev/swis/v5)
[![swis-api CI/CD pipeline](https://github.com/savla-dev/swis-api/actions/workflows/deployment.yml/badge.svg?branch=master)](https://github.com/savla-dev/swis-api/actions/workflows/deployment.yml)

sakalWebIS v5 RESTful JSON API in Go for Docker

`swapi` is a condensated data structure tree made into an Information System (IS) environment. It works as an independent in-memory database/cache for generic system components like `users`, `roles`, `projects`, `business_entities` and many more. Primarily, it is meant to be used in Docker as all components are written to work out of the box when there (no go runtime installment needed when run in Docker).

Read more in a short article on `swapi`:

+ https://krusty.space/projects/swis-api

### simple way to run (docker stack without a reverse-proxy)

```
git clone https://github.com/savla-dev/swis-api
cp .env.example .env
vi .env

make build run
```

or use Github Action Runner (self-hosted) to (re)delpoy `swapi` using Actions Secrets (Settings -> Secrets and variables -> Actions -> New repository secret), use these keys at least:
```
# without http/https prefix amd HTTP path
APP_URL 
# token used for redeployment dump and data reimport (`dump_data` and `import_data` workflow jobs)
ROOT_TOKEN

# if you want to log to Loki, define an URL (full) to the loki instance
LOKI_URL
# defaults to 'debug', specify 'release' in .env or in GHA secret for production
GIN_MODE
```

### install latest binary (go runtime installed required)

```
go install go.savla.dev/swis/v5@latest

GIN_MODE=release ROOT_TOKEN=xxx SERVER_PORT=8050 ~/go/bin/swis
```

`swapi` could be run as a single binary too. However, some environment constants have to be set when running it solitarily --- mainly `SERVER_PORT` and `ROOT_TOKEN` environment variables are required for the app's smooth start-up. 

### development

`swapi` development environment is made in Docker too, it is simple to run a local instance using:

```
make fmt dev
```

Then, the local instance could be accessed via http://localhost:8051. To contribute, simply make a new branch and create a new pull request.
