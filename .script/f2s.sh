#!/bin/bash

# f2s.sh
# facter to swapi data format parser
# outputs raw JSON
# krusty@savla.dev / Jan 30, 2024

die() {
	echo $@
	exit 1
}

FACTS_FACTER=(
	"is_virtual"
	"kernelversion"

	"memory.system.total_bytes"
	"memory.system.used_bytes"

	"networking.domain"
	"networking.hostname"
	"networking.fqdn"
	"networking.ip"
	"networking.network"

	"os.architecture"
	"os.family"
	"os.selinux.enabled"

	"processors.cores"

	"system_uptime.seconds"
	"timezone"
)

FACTS_SWAPI=(
	"is_virtual"
	"kernel_version"

	"memory_total_bytes"
	"memory_used_bytes"

	"net_domain"
	"net_hostname"
	"net_fqdn"
	"net_primary_ip"
	"net_primary_network"

	"os_arch"
	"os_family"
	"os_selinux_enabled"

	"proc_cores"

	"system_uptime_sec"
	"timezone"
)

# run checks to ensure used tools persistence
which facter > /dev/null || die "[ ! ] facter tool not found"
which jq > /dev/null || die "[ ! ] jq tool not found"

i=0

echo "{"

# loop over all facts defined above
for fact in "${FACTS_FACTER[@]}"; do
	# show the actual fact resolution
	#facter -j $fact 

	# save parsed value to a fragment var
	fragment=$(facter -j $fact | jq ".[\"$fact\"]")

	echo "\"${FACTS_SWAPI[$i]}\": $fragment,"

	# increment indexing variable
	i=$((i+1))
done

echo "\"timestamp\": $(date +%s)"
echo "}" 

