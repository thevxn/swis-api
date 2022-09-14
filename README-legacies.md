# swis-api v5 IS REST API legacies

# importing legacy

These manual commands can be batch executed using Makefile:

```
# dump prod data
make dump

# reimport data to prod (override current data set)
make import_prod_static_data
```

```shell
# load environment variables/constants
. .env

# (manual) import depot items example
curl -sLX POST \
	-H "X-Auth-Token: ${ROOT_TOKEN}" \
	-d @${BACKUP_DIR}/depots.json \
	http://localhost:${DOCKER_EXTERNAL_PORT}/depots/restore | jq .

# (manual) import users example
curl -sLX POST \
	-H "X-Auth-Token: ${ROOT_TOKEN}" \
	-d @${BACKUP_DIR}/users.json \
	http://localhost:${DOCKER_EXTERNAL_PORT}/users/restore | jq .

# (manual) import alvax command list example
curl -sLX POST \
	-H "X-Auth-Token: ${ROOT_TOKEN}" \
	-d @${BACKUP_DIR}/alvax_command_list.json \ 
	http://localhost:${DOCKER_EXTERNAL_PORT}/alvax/commands/restore | jq .

# (manual) import krusty's SSH keys example
curl -sLX POST \
	-H "X-Auth-Token: ${ROOT_TOKEN}" \
	-d @.dumps/krusty_ssh_keys.json \
	http://localhost:${DOCKER_EXTERNAL_PORT}/users/krusty/keys/ssh | jq .
```


### legacy MariaDB export (n0p_depot)

MySQL export + regexp dump parsing "case-study".

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


