#
# swis-api (swapi) / Makefile
#

#
# VARS
#

# load example variables/constants, override by .env file if exists
include .env.example
-include .env

APP_ENVIRONMENT?=development
PROJECT_NAME?=${APP_NAME}
DOCKER_COMPOSE_FILE?=deployments/docker-compose.yml
DOCKER_COMPOSE_OVERRIDE?=deployments/docker-compose.override.yml
DOCKER_COMPOSE_DEV_FILE?=deployments/docker-compose.dev.yml
DOCKER_COMPOSE_DEV_OVERRIDE?=deployments/docker-compose.dev.override.yml
SWAG_BINARY?=~/go/bin/swag

APP_URL?=swapi.example.com
LOKI_URL?=loki.example.com/loki/api/v1/push
ROOT_TOKEN?=${ROOT_TOKEN_DEFAULT}
GIN_MODE?=debug
#CF_API_EMAIL?=
#CF_API_TOKEN?=
#CF_BEARER_TOKEN?=

GOARCH := $(shell go env GOARCH)
GOCACHE?=/home/${USER}/.cache/go-build
GOMODCACHE?=/home/${USER}/go/pkg/mod
GOOS := $(shell go env GOOS)

PATH:=${PATH}:/usr/bin

# test env
POSTMAN_COLLECTION_FILE=test/postman/swapi_E2E_dish.postman_collection.json
HOSTNAME?=localhost
ROOT_TOKEN_TEST=fsdFD33FdsfK3dcc00Wef223kffDSrrrr

# define standard colors
# https://gist.github.com/rsperl/d2dfe88a520968fbc1f49db0a29345b9
ifneq (,$(findstring xterm,${TERM}))
	BLACK        := $(shell tput -Txterm setaf 0)
	RED          := $(shell tput -Txterm setaf 1)
	GREEN        := $(shell tput -Txterm setaf 2)
	YELLOW       := $(shell tput -Txterm setaf 3)
	LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
	PURPLE       := $(shell tput -Txterm setaf 5)
	BLUE         := $(shell tput -Txterm setaf 6)
	WHITE        := $(shell tput -Txterm setaf 7)
	RESET        := $(shell tput -Txterm sgr0)
else
	BLACK        := ""
	RED          := ""
	GREEN        := ""
	YELLOW       := ""
	LIGHTPURPLE  := ""
	PURPLE       := ""
	BLUE         := ""
	WHITE        := ""
	RESET        := ""
endif

export


#
# TARGETS
#

.PHONY: info
info: 
	@echo -e "\n${GREEN} ${PROJECT_NAME} / Makefile ${RESET}\n"

	@echo -e "${YELLOW} make${RESET} --- ${BLUE}show this helper ${RESET}\n"

	@echo -e "${YELLOW} make fmt${RESET}  --- ${BLUE}reformat the go source (gofmt) ${RESET}"
	@echo -e "${YELLOW} make docs${RESET} --- ${BLUE}render documentation from code (swagger OA docs) ${RESET}\n"

	@echo -e "${YELLOW} make build${RESET} --- ${BLUE}build project (docker image) ${RESET}"
	@echo -e "${YELLOW} make run${RESET}   --- ${BLUE}run project ${RESET}"
	@echo -e "${YELLOW} make logs${RESET}  --- ${BLUE}fetch container's logs ${RESET}"
	@echo -e "${YELLOW} make stop${RESET}  --- ${BLUE}stop and purge project (only docker containers!) ${RESET}\n"

	@echo -e "${YELLOW} make import_dump${RESET} --- ${BLUE}import dumped data (locally) into runtime ${RESET}"
	@echo -e "${YELLOW} make dump${RESET}        --- ${BLUE}dump runtime data to DUMP_DIR ${RESET}"
	@echo -e "${YELLOW} make backup${RESET}      --- ${BLUE}execute data dump and tar/gzip data backup ${RESET}"
	@echo -e ""

.PHONY: config
config: 
	@echo -e "\n${YELLOW} Installing additional tools and packages... ${RESET}\n"
	@go install github.com/swaggo/swag/cmd/swag@latest
	

.PHONY: version
version: 
	@echo -e "\n${YELLOW} Updating app's version (docs) according to dot-env file... ${RESET}\n"
	@/usr/bin/sed -i 's/\(\/\/[ ]@version\) .*/\1 ${APP_VERSION}/' cmd/swis-api/main.go

.PHONY: fmt
fmt:
	@echo -e "\n${YELLOW} Code reformating (gofmt)... ${RESET}\n"
	@/usr/local/go/bin/gofmt -s -w .
#@find . -name "*.go" -exec gofmt {} \;

.PHONY: build
build:  version
	@echo -e "\n${YELLOW} Building project (docker compose build)... ${RESET}\n"
	@docker compose --file $(DOCKER_COMPOSE_FILE) build
#@docker compose --file $(DOCKER_COMPOSE_FILE) build --no-cache

.PHONY: test
test:
	@echo -e "\n${YELLOW} Running tests in all packages (go test)... ${RESET}\n"
	go test -count=1 -v ./...

.PHONY: unit
unit: test

.PHONY: bench
bench:
	@echo -e "\n${YELLOW} Running benchmark tests (go test)... ${RESET}\n"
	go test -bench=. ./...

.PHONY: coverage
coverage:
	@echo -e "\n${YELLOW} Running code coverage evaluation (go test)... ${RESET}\n"
	go test -coverprofile -v ./... 
	go tool cover -html=coverage.out

.PHONY: test_deploy
test_deploy:
	@echo -e "\n${YELLOW} Starting temporary test container... ${RESET}\n"
	@docker run --rm --detach \
		--name ${DOCKER_TEST_CONTAINER_NAME} \
		-p ${DOCKER_TEST_PORT}:${DOCKER_TEST_PORT} \
		-e ROOT_TOKEN=${ROOT_TOKEN_TEST} \
		-e SERVER_PORT=${DOCKER_TEST_PORT} \
		${DOCKER_IMAGE_TAG}

.PHONY: e2e
e2e:	
	@echo -e "\n${YELLOW} Running Postman collection... ${RESET}\n"
	@postman collection run ${POSTMAN_COLLECTION_FILE} --timeout-request 5000 --env-var "token=${ROOT_TOKEN_TEST}" --env-var "baseUrl=${HOSTNAME}:${DOCKER_TEST_PORT}"; \
		docker stop ${DOCKER_TEST_CONTAINER_NAME}

APP_URLS_TRAEFIK?=`${APP_URL}`
.PHONY: run
run:
	@echo -e "\n${YELLOW} Starting project (docker compose up)... ${RESET}\n"
	@[ -f "${DOCKER_COMPOSE_DEV_OVERRIDE}"  ] \
		&& docker compose --file $(DOCKER_COMPOSE_FILE) --file ${DOCKER_COMPOSE_OVERRIDE} up --force-recreate --remove-orphans --detach \
		|| docker compose --file $(DOCKER_COMPOSE_FILE) up --force-recreate --remove-orphans --detach

.PHONY: logs
logs:
	@echo -e "\n${YELLOW} Fetching container's logs (CTRL-C to exit)... ${RESET}\n"
	@docker logs ${DOCKER_CONTAINER_NAME} --follow

.PHONY: stop
stop:  
	@echo -e "\n${YELLOW} Stopping and purging project (docker compose down)... ${RESET}\n"
	@docker compose --file $(DOCKER_COMPOSE_FILE) down

.PHONY: dev
dev: fmt
	@echo -e "\n${YELLOW} Starting local swapi instance... ${RESET}\n"
	@[ -f "${DOCKER_COMPOSE_DEV_OVERRIDE}"  ] \
		&& /usr/bin/docker compose --file $(DOCKER_COMPOSE_DEV_FILE) --file ${DOCKER_COMPOSE_DEV_OVERRIDE} up --force-recreate --remove-orphans --build \
		|| /usr/bin/docker compose --file $(DOCKER_COMPOSE_DEV_FILE) up --force-recreate --remove-orphans --build

.PHONY: dump
dump: 
	@echo -e "\n${YELLOW} Dumping prod data to ${DUMP_DIR}... ${RESET}\n"
	@scripts/dump_prod_data.sh

.PHONY: backup
backup: dump
	@echo -e "\n${YELLOW} Archiving and compressing dumped data... ${RESET}\n"
	@scripts/backup_dumped_files.sh

.PHONY: import_dump
import_dump: 
	@echo -e "\n${YELLOW} Import stored data (${DUMP_DIR}) to backend... ${RESET}\n"
	@scripts/import_prod_data.sh

.PHONY: push
push:
	@echo -e "\n${YELLOW} (re)tagging project and pushing to $(git branch --show-current) branch... ${RESET}\n"
	@/usr/bin/git tag -fa v${APP_VERSION} -m "v${APP_VERSION}"
	@/usr/bin/git push --set-upstream origin $(git branch --show-current) --follow-tags

.PHONY: docs
docs:
	@echo -e "\n${YELLOW} Regenerating documentation for swagger and rebuilding binary file... ${RESET}\n"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@${SWAG_BINARY} init --parseDependency -ot json . -g ./cmd/swis-api/main.go 
	@docker compose --file $(DOCKER_COMPOSE_FILE) up swagger_ui --detach --force-recreate

.PHONY: sh
sh: 
	@echo -e "\n${YELLOW} Attaching container's (${DOCKER_CONTAINER_NAME}) shell... ${RESET}\n"
	@docker exec -it ${DOCKER_CONTAINER_NAME} sh

USER_TOKEN?=xxx
TARGET_INSTANCE_URL?=http://localhost:${DOCKER_EXTERNAL_PORT}
URL_PATH?=/
METHOD?=GET
FLAGS?=-sL
.PHONY: raw
raw:
#@echo -e "\n${YELLOW} Executing a raw cURL request based on .env variables... ${RESET}\n"
	@/usr/bin/curl ${FLAGS} -X ${METHOD} -H "X-Auth-Token: ${USER_TOKEN}" ${TARGET_INSTANCE_URL}${URL_PATH}
	
.PHONY: fetch_facts
fetch_facts:
	@echo -e "\n${YELLOW} Executing a raw cURL request based on .env variables... ${RESET}\n"
	@scripts/fetch_facts.sh

.PHONY: compose_host_vars
compose_host_vars:
	@echo -e "\n${YELLOW} Fetching hosts configuration, exporting as YAML host_vars... ${RESET}\n"
	@scripts/compose_host_vars.sh
