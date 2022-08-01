#!/bin/sh

# import_prod_data.sh
# batch import static data from .data import/POST to prod.swapi.savla.su -> swapi.savla.su
# "workflow" chain structure on the bottom of this script -- below function blocks
# May 27, 2022 / krusty@savla.dev


# nonzero exit macro
function die {
	echo $@
	exit 1
}


#
# vars
#

# needs (but has defaults):
# APP_URL -> DEST_URL
# APP_ROOT -> DATA_DIR

DEST_URL="${APP_URL:-http://swapi.savla.su}"
DATA_DIR="${APP_ROOT:-./.data}"
[ ! -d "${DATA_DIR}" ] && die "DATA_DIR (${DATA_DIR)} of a no existence"

# tools test
[ ! -f "$(which curl)" ] && die "'curl' tool not found on runtime"
[ ! -f "$(which jq)" ] && die "'jq' tool not found on runtime"

# use explicitly POST method (-X), hide conn progress info (-s), follow locations (-L)
alias curlp="$(which curl) -v -sLX POST"
alias jq="$(which jq)"

#
# import blocks (to differenciate the script/workflow better...)
# should be self-documenting
#

function import_users {
  # import users, template:
  #curl -d @.data/users.json -sLX POST http://${APP_URL}/users/restore | jq .
  POST_PATH="/users/restore"
  URL="${DEST_URL}${POST_PATH}"

  DATA_FILE="${DATA_DIR}/users.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing users..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}

function import_ssh_keys {
  # template:
  #curl -d @.data/ssh_keys_tack.json -sLX POST http://swapi.savla.su/users/tack/keys/ssh | jq .

  # users to import given SSH keys arrays to (according to ./data/ssh_keys_username.json convention and swis-api/users.User.SSHKeys model)
  SSH_KEYS_USERS=(
    krusty
    stepis
    tack
  )

  # loop over users and import ssh keys for them
  for USER in ${SSH_KEYS_USERS[@]}; do
    POST_PATH="/users/${USER}/keys/ssh"
    URL="${DEST_URL}${POST_PATH}"

    DATA_FILE=${DATA_DIR}/ssh_keys_${USER}.json
    [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

    echo "imporitng SSH keys to ${USER}..."
    curlp --data @${DATA_FILE} --url ${URL} | jq .
  done
}

function import_depots {
  # import depot items; template:
  #curl -d @.data/depots.json -sLX POST http://${APP_URL}/depots/restore | jq .
  POST_PATH="/depots/restore"
  URL="${DEST_URL}${POST_PATH}"

  DATA_FILE="${DATA_DIR}/depots.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing/restoring/rewriting depots..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}

function import_alvax_cmd_list {
  # import alvax command list
  #curl -d @.data/alvax_command_list.json -sLX POST http://${APP_URL}/alvax/commands/restore | jq .
  POST_PATH="/alvax/commands/restore"
  URL="${DEST_URL}${POST_PATH}"

  DATA_FILE="${DATA_DIR}/alvax_command_list.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing alvax command list..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}

function import_dish_sockets {
  POST_PATH="/dish/sockets/restore"
  URL="${DEST_URL}${POST_PATH}"

  DATA_FILE="${DATA_DIR}/dish_sockets.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing dish sockets (socket list)..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}

function import_infrastructure {
  POST_PATH="/infra/restore"
  URL="${DEST_URL}${POST_PATH}"

  DATA_FILE="${DATA_DIR}/infra.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing infrastructure (hosts+networks)..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}

function import_finance {
  POST_PATH="/finance/restore"
  URL="${DEST_URL}${POST_PATH}"

  DATA_FILE="${DATA_DIR}/finance_accounts.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing finance (all accounts)..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}

function import_business_array {
  POST_PATH="/business/restore"
  URL="${DEST_URL}${POST_PATH}"

  DATA_FILE="${DATA_DIR}/business_array.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing business array..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}

function import_projects {
  POST_PATH="/projects/restore"
  URL="${DEST_URL}${POST_PATH}"

  DATA_FILE="${DATA_DIR}/projects.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing projects..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}


#
# workflow/chain import
#

import_users          || die "cannot import users"
import_ssh_keys       || die "cannot import SSH keys"
import_depots         || die "camnot import depots"
import_alvax_cmd_list || die "cannot import alvax command list"
import_dish_sockets   || die "cannot import dish sockets"
import_infrastructure || die "cannot import infrastructure"
import_finance 	      || die "cannot import finance"
import_business_array || die "cannot import business array"
import_projects	      || die "cannot import projects"

