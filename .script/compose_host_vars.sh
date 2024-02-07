#!/bin/bash

# compose_host_vars.sh
#
# krusty@savla.dev / Feb 7, 2024

die() {
	echo $@
	exit 1
}

[ -z "${USER_TOKEN}" -o -z "${TARGET_INSTANCE_URL}" ] && die "USER_TOKEN and/or TARGET_INSTANCE_URL env vars are empty"

mkdir -p ~/host_vars

# fetch hosts list
hosts=$(curl -sL -H "X-Auth-Token: ${USER_TOKEN}" -X GET "${TARGET_INSTANCE_URL}/infra/hosts" | jq -r ".items[] | .id")

# indexing helper
i=0

# iterate over a host, ensure deps installed, fetch the exported facts file and upload it to swapi
for host in ${hosts}; do
	echo "--- ${host}"
	cat <<-EOF > ~/host_vars/${host}
	---

	#
	# ${host}
	#

	EOF
	curl -sL -H "X-Auth-Token: ${USER_TOKEN}" -X GET "${TARGET_INSTANCE_URL}/infra/hosts/${host}" | jq -r .item.configuration | \
		yq -y >> ~/host_vars/${host}
	echo
done
