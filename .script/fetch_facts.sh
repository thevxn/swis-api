#!/bin/bash

# fetch_facts.sh
#
# krusty@savla.dev / Feb 6. 2024

die() {
	echo $@
	exit 1
}

# run the control for mandatory ones
[ -z "${USER_TOKEN}" -o -z ${TARGET_INSTANCE_URL} ] && die "no user token entered"

mkdir -p ~/cloud_facts

# fetch hosts list
hosts=$(curl -sL -H "X-Auth-Token: ${USER_TOKEN}" -X GET "${TARGET_INSTANCE_URL}/infra/hosts" | jq -r ".items[] | .id +\":\"+ .hostname_short")

# indexing helper
i=0

# iterate over a host, ensure deps installed, fetch the exported facts file and upload it to swapi
for host in ${hosts}; do
	id=$(echo $host | cut -d':' -f1)
	fqdn=$(echo $host | cut -d':' -f2)
	echo "--- ${fqdn} (${id}): start"

	pmn=$(ssh -A root@${fqdn} "dnf --version > /dev/null && echo dnf || echo apt")
	file=~/cloud_facts/${id}.json

	# copy f2s.sh to host, run it remotely, and fetch the output back
	scp .script/f2s.sh root@${fqdn}:
	ssh -A root@${fqdn} "${pmn} install -y facter jq && /root/f2s.sh > /root/${id}.json"
	scp root@${fqdn}:/root/${id}.json ${file}

	# upload facts to swapi
	ls ${file} || die "exported facts file not found locally"

	curl -sL -H "X-Auth-Token: ${USER_TOKEN}" -X POST --data @${file} "${TARGET_INSTANCE_URL}/infra/hosts/${id}/facts"| jq . \
		&& echo "--- ${fqdn} (${id}): facts uploaded" \
		|| die "error occured during the upload"

	i=$((i+1))
done
