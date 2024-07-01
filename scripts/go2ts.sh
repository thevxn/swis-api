#!/bin/bash

# go2ts.sh 
# helper script to regenerate TypeScript code from various pkgs' models
# krusty@savla.dev / July 1, 2024

die() {
	echo "[ fail ] $@"
	exit 1
}

warn() {
	echo "[ warn ] $@"
}

declare -a PKGS=(
	"alvax"
	"backups"
	"business"
	"depots"
	"dish"
	"finance"
	"infra"
	"links"
	"news"
	"projects"
	"queue"
	"roles"
	"users"
)

GO2TS_REPO=./third_party/go2ts
GO2TS_BIN=${GO2TS_REPO}/bin/go2ts

PKG_DIR=./pkg

TS_SRC_DIR=./third_party/swapi-types/src
TS_INDEX=${TS_SRC_DIR}/index.ts

# ensure 
[ -d "${GO2TS_REPO}" -a -d "${TS_SRC_DIR}" ] || {
	warn "submodules prolly not fetched, fetching now..."
	git submodule update --init --recursive --remote || die "unable to update submodules"
}

# check the existence of go2ts executable binary
[ -x "${GO2TS_BIN}" ] || {
	warn "go2ts binary not built, building now..."
	mkdir -p ${GO2TS_REPO}/bin
	go build -C ${GO2TS_REPO} -o bin/go2ts cmd/go2ts/main.go || die "unable to build the go2ts binary"
}

# trash out the contents of index.ts file
truncate --size 0 $TS_INDEX

# iterate over PKGS and generate .ts files from each item
for (( i=0; i<${#PKGS[@]}; i++ )); do
	PKG=${PKGS[$i]}
	echo ${PKG}

	# check if such item exists
	[ -d "${PKG_DIR}/${PKG}" ] || continue

	# run the GO2TS code translation
	${GO2TS_BIN} ${PKG_DIR}/${PKG} > ${TS_SRC_DIR}/${PKG}.ts 2> /dev/null && {
		echo "export * from './${PKG}'" >> ${TS_INDEX}
		continue
	} || {
		warn "pkg '${PKG}': translation failed, hotfixing..."
		mkdir -p ./tmp/${PKG}
		cat ${PKG_DIR}/${PKG}/models.go | sed -e '/chan/d' > ./tmp/${PKG}/models.go
	}

	# try again (hotfix for unsupported types.Chan issue)
	${GO2TS_BIN} ./tmp/${PKG} > ${TS_SRC_DIR}/${PKG}.ts && {
		warn "pkg '${PKG}': code translation hotfixed (chan type removed)..."
		echo "export * from './${PKG}'" >> ${TS_INDEX}
	} || {
		warn "pkg '${PKG}': code translation hotfix failed, skipping..."
	};

	rm -rf ./tmp/${PKG}
done


