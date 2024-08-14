#!/bin/bash

# migrate.sh
# simple helper script to migrate data from v5.16 to v5.17
# krusty@vxn.dev / Aug 14, 2024

[ -d "${DUMP_DIR}" ] || exit 1
cd ${DUMP_DIR}

jq . alvax_configs.json | \
	sed -e 's/"key": "\(.*\)",/"id": "\1",\n"key": "\1",/' | \
	jq . > alvax_configs.pretty.json && \
	mv alvax_configs.pretty.json alvax_configs.json

jq . backups.json | \
	sed -e 's/service_name/id/' | \
	sed -e 's/"backup_size":[[:space:]]"\([0-9]*\)",/"backup_size": \1/' > backups.pretty.json && \
	mv backups.pretty.json backups.json

jq . depots.json | \
	sed -e 's/"id":[[:space:]]\([0-9]*\),/"id": "\1"/' | \
	jq . > depots.pretty.json && \
	mv depots.pretty.json depots.json

jq . dish_sockets.json | \
	sed -e 's/socket_id/id/' > dish_sockets.pretty.json && \
	mv dish_sockets.pretty.json dish_sockets.json

jq . finance.json | \
	sed -e 's/account_id/id/' > finance.pretty.json && \
	mv finance.pretty.json finance.json

jq . infra.json | \
	sed -e 's/domain_id/id/' | \
	sed -e 's/"hash": "\(.*\)",/"id": "\1",\n"hash": "\1",/' | \
	jq . > infra.pretty.json && \
	mv infra.pretty.json infra.json

jq . links.json | \
	sed -e 's/"name": "\(.*\)",/"id": "\1",\n"name": "\1",/' | \
	jq . > links.pretty.json && \
	mv links.pretty.json links.json

jq . news_sources.json | \
	sed -e 's/source_id/id/' > news_sources.pretty.json && \
	mv news_sources.pretty.json news_sources.json

jq . projects.json | \
	sed -e 's/project_id/id/' > projects.pretty.json && \
	mv projects.pretty.json projects.json

jq . roles.json | \
	sed -e 's/"name": "\(.*\)",/"id": "\1",\n"name": "\1",/' | \
	jq . > roles.pretty.json && \
	mv roles.pretty.json roles.json

jq . users.json | \
	sed -e 's/"name": "\(.*\)",/"id": "\1",\n"name": "\1",/' | \
	jq . > users.pretty.json && \
	mv users.pretty.json users.json

