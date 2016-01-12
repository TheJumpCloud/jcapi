#!/bin/bash
set -e

BUILD_PATH="`pwd`/build"
mkdir -p $BUILD_PATH

function go_build () {
  OUTPUT_PREFIX=${PWD##*/}
  GOOS="$1"
  GOARCH="$2"
  go build -o "${BUILD_PATH}/${OUTPUT_PREFIX}_${GOOS}_${GOARCH}"
}

function cross_build () {
  go_build "darwin" "amd64"
  go_build "linux" "386"
  go_build "linux" "amd64"
  go_build "windows" "386"
  go_build "windows" "amd64"
}

function build_directory () {
  echo $@
  pushd $1
  cross_build
  popd
}

FOLDERS_TO_BUILD=(
  "examples/ExportUsersPerSystemToCSV/"
)

for path in "${FOLDERS_TO_BUILD[@]}"; do
  build_directory $path
done

