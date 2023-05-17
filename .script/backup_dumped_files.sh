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
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
STATUS=failure
SIZE=""

mkdir -p ${BACKUP_TARGET_DIR}

tar --exclude='*.tar.gz' -czvf ${BACKUP_TARGET_DIR}/${TIMESTAMP}.tar.gz ${BACKUP_TARGET_DIR}/.. 
[[ $? -ne 0 ]] && die "backup error" || {
	SIZE=$(du -shx ${BACKUP_TARGET_DIR}/${TIMESTAMP}.tar.gz | awk '{ print $1 }')
	STATUS=success
}

# report back to swapi/backups
URL=http://localhost:${DOCKER_EXTERNAL_PORT}
SERVICE_NAME=swis-api
TOKEN=xxxDS3RKKSddK43KDLSA34AAa4AAAA

# generate backup report
TMP_FILE_NAME=/tmp/swis-backup-report.json
cat > ${TMP_FILE_NAME} <<-EOF
{
        "service_name": "${SERVICE_NAME}",
        "last_status":  "${STATUS}",
        "backup_size":  "${SIZE}",
        "timestamp":    $(date +%s),
        "file_name":    "${TIMESTAMP}.tar.gz"
}
EOF

curl -X PUT -sL -H "X-Auth-Token: $TOKEN" \
        --data @${TMP_FILE_NAME} \
        ${URL}/backups/${SERVICE_NAME}

echo "backup successful --- ${BACKUP_TARGET_DIR}/${TIMESTAMP}.tar.gz"
