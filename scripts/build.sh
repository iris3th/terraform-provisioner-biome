#!/bin/bash

set -e
# Prerequisites
if ! command -v gox > /dev/null; then
  go get -u github.com/mitchellh/gox

fi

# setup environment
PROVISIONER_NAME="biome_dev"
TARGET_DIR="$(pwd)/results"
XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
XC_OS=${XC_OS:=linux darwin windows freebsd openbsd solaris}
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386 !linux/arm !linux/386 !freebsd/386 !freebsd/amd64 !freebsd/arm !openbsd/386 !openbsd/amd64 !solaris/amd64"
LD_FLAGS="-s -w"
export CGO_ENABLED=0

rm -rf "${TARGET_DIR}"
mkdir -p "${TARGET_DIR}"

# Compile
gox \
  -os="${XC_OS}" \
  -arch="${XC_ARCH}" \
  -osarch="${XC_EXCLUDE_OSARCH}" \
  -ldflags "${LD_FLAGS}" \
  -output "$TARGET_DIR/{{.OS}}_{{.Arch}}/terraform-provisioner-${PROVISIONER_NAME}_v0.1" \
  -verbose \
  -rebuild \
