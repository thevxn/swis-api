# swis-api (swapi)

[![swis-api CI/CD pipeline](https://github.com/savla-dev/swis-api/actions/workflows/docker-image.yml/badge.svg)](https://github.com/savla-dev/swis-api/actions/workflows/docker-image.yml)

[sakalWeb (v5) Information System RESTful API (intranet)](http://swapi.savla.su)

## repo vademecum

### .assets

This folder contains static assets like images (favicon), that are served by server statically (not dynamically).

### .bin

This folder contains executables for swapi package management, like data importing batch script.

### .data

Crutial folder for all data to be imported to a running instance, for this purpose, .bin/ importing scriped is used. This folder is used as temporary data backup.

### .env

Capital configuration file containing environment constants for `Makefile` and other linked files like `docker-compose.yml`.

### main.go

Main package's logic is coded there, main package's entrypoint.

### modules

All swapi modules are stored in their folders. Every module has its `models.go` file with data structures, `controllers.go` with its methods and functions, and `routes.go` for gin router to serve module's handles.

### Dockerfile

Recipe for docker image build using `docker build .`

### Makefile

Dev/build recipe for GNU `make` tool. Listing (Jan 14, 2022):

```shell
$ make

 swis-core-api / Makefile 

 make build   --- build project (docker image) 
 make run     --- run project 
 make log     --- fetch container's log 
 make stop    --- stop and purge project (only docker containers!) 

```

### docker-compose.yml

YAML-formated file for docker-compose stack. Contains defitions for docker container and its isolated network. Uses constants from `.env` dot-file.


## documentation

[swagger 2.0 is used to document API scheme](http://swapi-docs.savla.su)

```
# at project root run
swag init .
make docs

# build local binary
go build swis-api

# run server
./swis-api

# view
http://localhost:8049/swagger/index.html
```

## authentication

Swapi uses token-based authentication for any request to be authenticated, For initial importing, `ROOT_TOKEN` (see [`.env`](/.env) file) is used by importing executable. For any request, the header `X-Auth-Token` has to be sent with a custom HTTP request.

## importing

At start, swapi instance memory is cleared and ready for any data import (until the next restart). Any data stored in runtime memory should be dumped using GET methods at particular paths. This approach should make `swapi` instance universal (while omitting custom packages/modules).

```shell
# run local instance (redeployment CI/CD job)
make build run

# import prod data -- local .data files to swapi.savla.su prod URL (import_data CI/CD job)
make import_prod_static_data
```

```shell
# (manual) import depot items example
curl -d @.data/depot.json -sLX POST http://localhost:8003/depot/restore | jq .

# (manual) import users example
curl -d @.data/users.json -sLX POST http://localhost:8003/users/restore | jq .

# (manual) import alvax command list example
curl -d @.data/alvax_command_list.json -sLX POST http://localhost:8003/alvax/commands/restore | jq .

# (manual) import SSH keys example
curl -d @.data/krusty_ssh_keys.json -sLX POST http://swapi.savla.su/users/krusty/keys/ssh
```


### legacy MariaDB export (n0p_depot)

```shell
# export legacy table contents, and reformat result lines into JSON array items
mysql -u n0p_sysadm -p n0p_core -sNe 'select JSON_ARRAY(id, n0p_depot.desc, misc, depot) from n0p_depot;' > n0p_depot.export.json

# check correctness of a JSON file (has to pass, ergo exitcode == 0)
jq . n0p_depot.export.json

# regexp for bracket change: [ -> {, ] -> }
2,$s/\[/\{/g
2,$s/\]/\}/g

# insert a comma ',' at the EOL
2,$s/^\(.*\)$/\1,/

# take all for array items and convert them into a JSON object
2,342s/^{\(.*\),[ ]\(".*"\),[ ]\(".*"\),[ ]\(".*"\)\},$/\{"id": \1, "desc": \2, "misc": \3, "depot": \4},/
```


