#!/bin/sh

# dump_prod_data.sh
# batch dump production data
# Sep 10, 2022 / krusty@savla.dev


# nonzero exit macro
function die {
	echo $@
	exit 1
}


#
# directory settings
#

DEST_URL="${APP_URL:-http://swapi.savla.su}"

[ -z ${DUMP_DIR} ] && die "DUMP_DIR constant unset!"
DATA_DIR=${DUMP_DIR}

mkdir -p ${DATA_DIR}


#
# tools test and reconfig
#

[ ! -f "$(which curl)" ] && die "'curl' tool not found on runtime"
[ ! -f "$(which jq)" ] && die "'jq' tool not found on runtime"


# use explicitly POST method (-X), hide conn progress info (-s), follow locations (-L)
alias curlp="$(which curl) -sLX GET -H 'X-Auth-Token: ${ROOT_TOKEN}'"
alias jq="$(which jq)"


function dump_generic {
  REQ_PATH="$1"
  URL="${DEST_URL}${REQ_PATH}"

  printf "dumping $2...\n\t"
  curlp --url ${URL} | tee ${DATA_DIR}/$2 | jq -r '. | {code,message} | join(" ")'
  echo
}

#
# modules
#

declare -a paths=(
	"/backups/"
	"/business/"
	"/depots/"
	"/dish/sockets"

	"/finance/"
	"/infra/"
	"/links/"
	"/news/sources"
	"/projects/"

	"/roles/"
	"/six/"
	"/users/"
)
declare -a files=(
	"/backups.json"
	"/business_array.json"
	"/depots.json"
	"/dish_sockets.json"

	"/finance_accounts.json"
	"/infra.json"
	"/links.json"
	"/news_sources.json"
	"/projects.json"

	"/roles.json"
	"/six_struct.json"
	"/users.json"
)

# restore all paths with files
for (( i=0; i<${#paths[@]}; i++ )); do
	dump_generic ${paths[$i]} ${files[$i]} || die "problem dumping ${files[$i]}"
done

