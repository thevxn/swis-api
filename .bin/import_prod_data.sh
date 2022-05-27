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

# use explicitly POST method (-X), hide conn progress info (-s), follow locations (-L)
alias curlp="$(which curl) -sLX POST"
alias jq="$(which jq)"


#
# import blocks (to differenciate the script/workflow better...)
# should be self-documenting
#

function import_users {
  # import users, template:
  #curl -d @.data/users.json -sLX POST http://${APP_URL}/users/restore | jq .
  _PATH="/users/restore"
  URL="${DEST_URL}${_PATH}"

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
    tack
  )

  # loop over users and import ssh keys for them
  for USER in ${SSH_KEYS_USERS[@]}; do
    _PATH="/users/${USER}/keys/ssh"
    URL="${DEST_URL}${_PATH}"

    DATA_FILE=${DATA_DIR}/ssh_keys_${USER}.json
    [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

    echo "imporitng SSH keys to ${USER}..."
    curlp --data @${DATA_FILE} --url ${URL} | jq .
  done
}

function import_depots {
  # import depot items; template:
  #curl -d @.data/depots.json -sLX POST http://${APP_URL}/depots/restore | jq .
  _PATH="/depots/restore"
  URL="${DEST_URL}${_PATH}"

  DATA_FILE="${DATA_DIR}/depots.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing/restoring/rewriting depots..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}

function import_alvax_cmd_list {
  # import alvax command list
  #curl -d @.data/alvax_command_list.json -sLX POST http://${APP_URL}/alvax/commands/restore | jq .
  _PATH="/alvax/commands/restore"
  URL="${DEST_URL}${_PATH}"

  DATA_FILE="${DATA_DIR}/alvax_command_list.json"
  [ ! -f "${DATA_FILE}" ] && die "DATA_FILE (${DATA_FILE}) of a no existence"

  echo "importing alvax command list..."
  curlp --data @${DATA_FILE} --url ${URL} | jq .
}


#
# workflow/chain import
#

import_users          || die "cannot import users"
import_ssh_keys       || die "cannot import SSH keys"
import_depots         || die "camnot import depots"
import_alvax_cmd_list || die "cannot import alvax command list"

