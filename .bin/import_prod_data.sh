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
#DATA_DIR="${APP_ROOT:-./.data}"
DATA_DIR=./.dumps
[ ! -d "${DATA_DIR}" ] && die "DATA_DIR (${DATA_DIR}) of a no existence"

# tools test
[ ! -f "$(which curl)" ] && die "'curl' tool not found on runtime"
[ ! -f "$(which jq)" ] && die "'jq' tool not found on runtime"

# use explicitly POST method (-X), hide conn progress info (-s), follow locations (-L)
alias curlp="$(which curl) -sLX POST -H 'X-Auth-Token: ${ROOT_TOKEN}'"
alias jq="$(which jq)"


#
# import blocks (to differenciate the script/workflow better...)
# should be self-documenting
#

function import_generic {
  POST_PATH="$1"
  URL="${DEST_URL}${POST_PATH}"

  DATA_FILE="${DATA_DIR}$2"
  [ ! -f "${DATA_FILE}" ] && die "[import_generic] DATA_FILE (${DATA_FILE}) of a no existence"

  printf "importing $2...\n\t"
  curlp --data @${DATA_FILE} --url ${URL} | jq -r '. | {code,message} | join(" ")'
  echo
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
    [ ! -f "${DATA_FILE}" ] && die "[import_ssh_keys] DATA_FILE (${DATA_FILE}) of a no existence"

    #echo "imporitng SSH keys to ${USER}..."
    curlp --data @${DATA_FILE} --url ${URL} | jq '. | {code,message} | join(" ")'
  done
}


#
# importing
#

declare -a paths=(
	"/alvax/commands/restore"
	"/backups/restore"
	"/business/restore"
	"/depots/restore"
	"/dish/sockets/restore"
	"/infra/restore"
	"/finance/restore"
	"/projects/restore"
	"/roles/restore"
	"/swife/restore"
	"/users/restore"
)
declare -a files=(
	"/alvax_command_list.json"
	"/backups.json"
	"/business_array.json"
	"/depots.json"
	"/dish_sockets.json"
	"/infra.json"
	"/finance_accounts.json"
	"/projects.json"
	"/roles.json"
	"/swife_frontends.json"
	"/users.json"
)

# restore all paths with files
for (( i=0; i<${#paths[@]}; i++ )); do
	import_generic ${paths[$i]} ${files[$i]} || die "problem importing ${files[$i]}"
done

# eventually import SSH keys too
import_ssh_keys

