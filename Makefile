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
DOCKER_COMPOSE_FILE?=./docker-compose.yml
SWAG_BINARY?=~/go/bin/swag

APP_URL?=swapi.example.com
LOKI_URL?=loki.example.com/loki/api/v1/push
ROOT_TOKEN?=${ROOT_TOKEN_DEFAULT}
APP_URLS_TRAEFIK?=`${APP_URL}`,`swis-api-run.local`,`swis-api.example.com`
GIN_MODE?=debug

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

	@echo -e "${YELLOW} make --- show this helper ${RESET}\n"

	@echo -e "${YELLOW} make fmt  --- reformat the go source (gofmt) ${RESET}"
	@echo -e "${YELLOW} make docs --- render documentation from code (swagger OA docs) ${RESET}\n"

	@echo -e "${YELLOW} make build --- build project (docker image) ${RESET}"
	@echo -e "${YELLOW} make run   --- run project ${RESET}"
	@echo -e "${YELLOW} make logs  --- fetch container's logs ${RESET}"
	@echo -e "${YELLOW} make stop  --- stop and purge project (only docker containers!) ${RESET}\n"

	@echo -e "${YELLOW} make import_dump --- import dumped data (locally) into runtime ${RESET}"
	@echo -e "${YELLOW} make dump        --- dump runtime data to DUMP_DIR ${RESET}"
	@echo -e "${YELLOW} make backup      --- execute data dump and tar/gzip data backup ${RESET}"
	@echo -e ""

.PHONY: version
version: 
	@echo -e "\n${YELLOW} Updating app's version (docs) according to dot-env file... ${RESET}\n"
	@sed -i 's/\(\/\/[ ]@version\) .*/\1 ${APP_VERSION}/' main.go

.PHONY: fmt
fmt:
	@echo -e "\n${YELLOW} Code reformating (gofmt)... ${RESET}\n"
	@gofmt -d -s .
#@find . -name "*.go" -exec gofmt {} \;

.PHONY: unit
unit:	
	@echo -e "\n${YELLOW} Running tests in all packages (go test)... ${RESET}\n"
	@go test -v ./...

.PHONY: build
build:  version
	@echo -e "\n${YELLOW} Building project (docker compose build)... ${RESET}\n"
	@docker compose --file $(DOCKER_COMPOSE_FILE) build
#@docker compose --file $(DOCKER_COMPOSE_FILE) build --no-cache

.PHONY: devploy
devploy:
	@echo -e "\n${YELLOW} Starting temporary dev container... ${RESET}\n"
	@docker run --rm --detach \
		--name ${DOCKER_DEV_CONTAINER} \
		-p ${DOCKER_DEV_PORT}:${DOCKER_DEV_PORT} \
		-e ROOT_TOKEN=d3qySD87Ds48300pl \
		-e DOCKER_INTERNAL_PORT=${DOCKER_DEV_PORT} \
		${DOCKER_IMAGE_TAG}

.PHONY: e2e
e2e:	
	@docker stop swis-api-dev

.PHONY: run
run:
	@echo -e "\n${YELLOW} Starting project (docker compose up)... ${RESET}\n"
	@docker compose --file $(DOCKER_COMPOSE_FILE) up --force-recreate --remove-orphans --detach

.PHONY: logs
logs:
	@echo -e "\n${YELLOW} Fetching container's logs (CTRL-C to exit)... ${RESET}\n"
	@docker logs ${DOCKER_CONTAINER_NAME} --follow

.PHONY: stop
stop:  
	@echo -e "\n${YELLOW} Stopping and purging project (docker compose down)... ${RESET}\n"
	@docker compose --file $(DOCKER_COMPOSE_FILE) down

.PHONY: dump
dump: 
	@echo -e "\n${YELLOW} Dumping prod data to ${DUMP_DIR}... ${RESET}\n"
	@.script/dump_prod_data.sh

.PHONY: backup
backup: dump
	@echo -e "\n${YELLOW} Archiving and compressing dumped data... ${RESET}\n"
	@.script/backup_dumped_files.sh

.PHONY: import_dump
import_dump: 
	@echo -e "\n${YELLOW} Import stored data (${DUMP_DIR}) to backend... ${RESET}\n"
	@.script/import_prod_data.sh

.PHONY: push
push:
	@echo -e "\n${YELLOW} (re)tagging project and pushing to master... ${RESET}\n"
	@git tag -fa v${APP_VERSION} -m "v${APP_VERSION}"
	@git push --follow-tags origin master

.PHONY: docs
docs:
	@echo -e "\n${YELLOW} Regenerating documentation for swagger and rebuilding binary file... ${RESET}\n"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@${SWAG_BINARY} init . 
	@docker compose --file $(DOCKER_COMPOSE_FILE) up swagger_ui --detach

.PHONY: sh
sh: 
	@echo -e "\n${YELLOW} Attaching container's (${DOCKER_CONTAINER_NAME}) shell... ${RESET}\n"
	@docker exec -it ${DOCKER_CONTAINER_NAME} sh
