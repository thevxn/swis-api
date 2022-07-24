#
# swis-core-api / Makefile
#

#
# VARS
#

-include .env

PROJECT_NAME?=swapi

DOCKER_DEV_IMAGE?=${PROJECT_NAME}-build
DOCKER_DEV_CONTAINER?=${PROJECT_NAME}-dev-run

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

# docker-compose vs docker compose (new syntax) check
COMPOSE_CMD:='docker compose'
ifeq (, $(shell which ${COMPOSE_CMD} 2>/dev/null))
	COMPOSE_CMD='docker-compose'
endif

export


#
# TARGETS
#

.PHONY: all info build go src make doc

all: info

info: 
	@echo -e "\n${GREEN} ${PROJECT_NAME} / Makefile ${RESET}\n"

#@echo -e "${YELLOW} make config  --- check dev environment ${RESET}"
	@echo -e "${YELLOW} make fmt     --- reformat the go source (gofmt) ${RESET}"
	@echo -e "${YELLOW} make doc     --- render documentation from code (go doc) ${RESET}\n"

	@echo -e "${YELLOW} make build   --- build project (docker image) ${RESET}"
	@echo -e "${YELLOW} make run     --- run project ${RESET}"
	@echo -e "${YELLOW} make logs    --- fetch container's logs ${RESET}"
	@echo -e "${YELLOW} make stop    --- stop and purge project (only docker containers!) ${RESET}"
	@echo -e ""

# target to see the runtime contents of COMPOSE_CMD constant -- to be deleted later
config:
	@echo ${COMPOSE_CMD}	

fmt:
	@echo -e "\n${YELLOW} Code reformating (gofmt)... ${RESET}\n"
	@gofmt -d -s .
	@find . -name "*.go" -exec gofmt {} \;

build: 
	@echo -e "\n${YELLOW} Building project (${COMPOSE_CMD} build)... ${RESET}\n"
	@$(COMPOSE_CMD) build --no-cache

run:	build
	@echo -e "\n${YELLOW} Starting project (${COMPOSE_CMD} up)... ${RESET}\n"
	@$(COMPOSE_CMD) up --force-recreate --detach

logs:
	@echo -e "\n${YELLOW} Fetching container's logs (CTRL-C to exit)... ${RESET}\n"
	@docker logs ${DOCKER_DEV_CONTAINER} -f

stop:  
	@echo -e "\n${YELLOW} Stopping and purging project (${COMPOSE_CMD} down)... ${RESET}\n"
	@$(COMPOSE_CMD) down

import_prod_static_data: 
	@echo -e "\n${YELLOW} Import stored data to backend... ${RESET}\n"
	@.bin/import_prod_data.sh

