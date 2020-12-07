#!/usr/bin/env bash
set -euo pipefail

$(dirname $0)/../scripts/build.sh

OS=$(docker info --format '{{.OSType}}')
pack package-buildpack dependency-cnb:${OS} --config package-${OS}.toml

if [[ "$OS" == "windows" ]]; then
BUILDER=cnbs/sample-builder:nanoserver-1809
TEST_COMMAND=(cmd /c "type c:\layers\my-dependency\my-layer\my-dependency\my-file.txt")
else
BUILDER=cnbs/sample-builder:bionic
TEST_COMMAND=(bash -c 'cat /layers/my-dependency/my-layer/my-dependency/my-file.txt')
fi

pack build dependency-app:${OS} \
  --buildpack docker://dependency-cnb:${OS} \
  --builder $BUILDER

docker run -i --rm dependency-app:${OS} -- "${TEST_COMMAND[@]}"
