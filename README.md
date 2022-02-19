# savla-dev/swis-api
sakalWeb IS RESTful API v5 [golang]

+ http://swapi.savla.su (intranet)

## repo files

### .env

dot-file containing base contants for Makefile, Dockerfile (via docker-compose.yml).

### Dockerfile

Recipe for docker image build using `docker build .`

### LICENSE

MIT license file for swis-core-api repository/project.

### Makefile

Dev/build recipe for GNU `make` tool. Listing (Jan 14, 2022):

```
$ make

 swis-core-api / Makefile 

 make build   --- build project (docker image) 
 make run     --- run project 
 make log     --- fetch container's log 
 make stop    --- stop and purge project (only docker containers!) 

```

### docker-compose.yml

YAML-formated file for docker-compose stack. Contains defitions for docker container and its isolated network. Uses constants from `.env` dot-file.

### go.mod (go.sum)

Application/package dependencies. Deprecated since the module is built inside the docker container.

### main.go

Source code file containing main() function.

