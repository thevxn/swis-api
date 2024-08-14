#!/bin/bash

# migrate.sh
# simple helper script to migrate data from v5.16 to v5.17
# krusty@vxn.dev / Aug 14, 2024

[ -d "${DUMP_DIR}" ] || exit 1
cd ${DUMP_DIR}

echo "Migrating alvax pkg..."
jq '.items[] | . += {"id": .key} | [.] | to_entries[] | {(.value.id): .value}' alvax_configs.json | jq -s 'add | {"items": .}' > alvax_configs.pretty.json || exit 1
mv alvax_configs.pretty.json alvax_configs.json

echo "Migrating backups pkg..."
jq '.items[] | with_entries(if .key == "service_name" then .key = "id" else . end) | .backup_size = (.backup_size | tonumber) | [.] | to_entries[] | {(.value.id): .value}' backups.json | jq -s 'add | {"items": .}' > backups.pretty.json || exit 1
mv backups.pretty.json backups.json
	
echo "Migrating depots pkg..."
jq '.items[] | . += {"id": (.id | tostring)} | [.] | to_entries[] | {(.value.id): .value}' depots.json | jq -s 'add | {"items": .}' > depots.pretty.json || exit 1
mv depots.pretty.json depots.json

echo "Migrating dish pkg..."
jq '.sockets[] | with_entries(if .key == "socket_id" then .key = "id" else . end) | [.] |to_entries[] | {(.value.id): .value}' dish_sockets.json | jq -s 'add | {"sockets": .}' > dish_sockets.pretty.json || exit 1
jq '.incidents | {"incidents": .}' dish_sockets.json > dish_incidents.pretty.json || exit 2
jq -s 'add' dish_sockets.pretty.json dish_incidents.pretty.json > dish_sockets.json

echo "Migrating finance pkg..."
jq '.items[] | with_entries(if .key == "account_id" then .key = "id" else . end) | [.] | to_entries[] | {(.value.id): .value}' finance.json | jq -s 'add | {"items": .}' > finance_accounts.pretty.json || exit 1
jq '.items | {"items": .}' finance.json > finance_items.pretty.json || exit 2
jq -s 'add' finance_accounts.pretty.json finance_items.pretty.json > finance.json
	
echo "Migrating infra pkg..."
jq '.domains[] | with_entries(if .key == "domain_id" then .key = "id" else . end) | [.] | to_entries[] | {(.value.id): .value}' infra.json | jq -s 'add | {"domains": .}' > infra_domains.pretty.json || exit 1
jq '.hosts | {"hosts": .}' infra.json > infra_hosts.pretty.json || exit 2
jq '.networks[] | . += {"id": .hash} | [.] | to_entries[] | {(.value.id): .value}' infra.json | jq -s 'add | {"networks": .}' > infra_networks.pretty.json || exit 3
jq -s 'add' infra_domains.pretty.json infra_hosts.pretty.json infra_networks.pretty.json > infra.json

echo "Migrating links pkg..."
jq '.items[] | . += {"id": .name} | [.] | to_entries[] | {(.value.id): .value}' links.json | jq -s 'add | {"items": .}' > links.pretty.json || exit 1
mv links.pretty.json links.json

echo "Migrating news pkg..."
jq '.items = (.items | (to_entries | map({(.key): {id: .key, user_name: .key, news_sources: .value}}) | add))' news_sources.json > news_sources.pretty.json || exit 1
mv news_sources.pretty.json news_sources.json

echo "Migrating projects pkg..."
jq '.items[] | with_entries(if .key == "project_id" then .key = "id" else . end) | [.] | to_entries[] | {(.value.id): .value}' projects.json | jq -s 'add | {"items": .}' > projects.pretty.json || exit 1
mv projects.pretty.json projects.json

echo "Migrating roles pkg..."
jq '.items[] | . += {"id": .name} | [.] | to_entries[] | {(.value.id): .value}' roles.json | jq -s 'add | {"items": .}' > roles.pretty.json || exit 1
mv roles.pretty.json roles.json

echo "Migrating users pkg..."
jq '.items[] | . += {"id": .name} | [.] | to_entries[] | {(.value.id): .value}' users.json | jq -s 'add | {"items": .}' > users.pretty.json || exit 1
mv users.pretty.json users.json

