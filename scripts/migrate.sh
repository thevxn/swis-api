#!/bin/bash

# migrate.sh
# simple helper script to migrate data from v5.16 to v5.17
# krusty@vxn.dev / Aug 14, 2024

[ -d "${DUMP_DIR}" ] || exit 1
cd ${DUMP_DIR}

echo "Migrating alvax pkg..."
jq '.items[] += {"id": .key}' alvax_configs.json > alvax_configs.pretty.json && \
	mv alvax_configs.pretty.json alvax_configs.json

echo "Migrating backups pkg..."
jq '.items[].id = .items[].service_name | del(.items[].service_name) | .items[].backup_size = (.items[].backup_size | tonumber)' backups.json > backups.pretty.json && \
	mv backups.pretty.json backups.json

echo "Migrating depots pkg..."
jq '.items[].id = (.items[].id | tostring)' depots.json > depots.pretty.json && \
	mv depots.pretty.json depots.json

echo "Migrating dish pkg..."
jq '.sockets[].id = .sockets[].socket_name | del(.sockets[].socket_id) ' dish_sockets.json > dish_sockets.pretty.json && \
	mv dish_sockets.pretty.json dish_sockets.json

echo "Migrating finance pkg..."
jq '.accounts[].id = .accounts[].account_id | del(.accounts[].account_id)' finance.json > finance.pretty.json && \
	mv finance.pretty.json finance.json

echo "Migrating infra pkg..."
jq '.domains[].id = .domains[].domain_id | del(.domains[].domain_id) | .networks[].id = .networks[].hash' infra.json > infra.pretty.json && \
	mv infra.pretty.json infra.json

echo "Migrating links pkg..."
jq '.items[].id = .items[].name' links.json > links.pretty.json && \
	mv links.pretty.json links.json

echo "Migrating news pkg..."
jq '.items = (.items | (to_entries | map({(.key): {id: .key, user_name: .key, news_sources: .value}}) | add))' news_sources.json > news_sources.pretty.json && \
	mv news_sources.pretty.json news_sources.json

echo "Migrating projects pkg..."
jq '.items[].id = .items[].project_id | del(.items[].project_name)' projects.json > projects.pretty.json && \
	mv projects.pretty.json projects.json

echo "Migrating roles pkg..."
jq '.items[].id = .items[].name' roles.json > roles.pretty.json && \
	mv roles.pretty.json roles.json

echo "Migrating users pkg..."
jq '.items[].id = .items[].name' users.json > users.pretty.json && \
	mv users.pretty.json users.json

