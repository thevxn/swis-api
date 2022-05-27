# swis-api (swapi)
sakalWeb IS RESTful API v5 [golang]

+ http://swapi.savla.su (intranet)

## importing

At start, swapi instance memory is cleared and ready for any data import (until the next restart). Any data stored in runtime memory should be dumped using GET methods at particular paths. This approach should make `swapi` instance universal (while omitting custom packages/modules).

```
# run local instance (redeployment CI/CD job)
make run

# import prod data -- local .data files to swapi.savla.su prod URL (import_data CI/CD job)
make import_prod_static_data


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

```
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

