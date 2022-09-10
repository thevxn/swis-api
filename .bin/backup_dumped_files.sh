#!/bin/sh

# backup_dumped_files.sh
# krusty@savla.dev / Sep 10, 2022

function die() {
	echo $@
	exit 1
}

# loaded by Makefile -- make backup to ensure having env const. loaded
[ -z ${DUMP_DIR} ] && die "DUMP_DIR not set"
BACKUP_TARGET_DIR=${DUMP_DIR}/archives
TIMESTAMP=$(date +"%d-%m-%Y_%H-%M-%S")
STATUS=failure
SIZE=""

mkdir -p ${BACKUP_TARGET_DIR}

tar --exclude='*.tar.gz' -czvf ${BACKUP_TARGET_DIR}/${TIMESTAMP}.tar.gz ${BACKUP_TARGET_DIR}/.. 

[[ $? -ne 0 ]] && die "backup error" || {
	SIZE=$(du -shx ${BACKUP_TARGET_DIR}/${TIMESTAMP}.tar.gz | awk '{ print $1 }')
	STATUS=success
}

# report back to swapi/backups
SERVICE_NAME=swis-api
TIMESTAMP=$(date +%s)
TOKEN=i1d67db229e54a3f73047d84cc4faea850c7ed38ae3c3069c9952731ffe20744ef798feda3c559075ebc2f3f350f5dedbc56ee8a7453537d401531c32881860b0

# generate backup report
cat > /tmp/backup-report.json <<-EOF
{
        "service_name": "${SERVICE_NAME}",
        "last_status":  "${STATUS}",
        "backup_size":  "${SIZE}",
        "timestamp":    ${TIMESTAMP},
        "file_name":    "${TIMESTAMP}.tar.gz"
}
EOF

curl -X PUT -sL -H "X-Auth-Token: $TOKEN" \
        --data @/tmp/backup-report.json \
        http://swapi.savla.su/backups/${SERVICE_NAME}

echo "backup successfull --- ${BACKUP_TARGET_DIR}/${TIMESTAMP}.tar.gz"
