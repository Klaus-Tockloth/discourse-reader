#!/bin/sh

# ------------------------------------
# Purpose:
# - Builds executables / binaries.
#
# Releases:
# - v1.0.0 - 2022-11-04: initial release
#
# Remarks:
# - go tool dist list
# ------------------------------------

# set -o xtrace
set -o verbose

# compile 'aix'
env GOOS=aix GOARCH=ppc64 go build -o binaries/aix-ppc64/discourse-reader

# compile 'darwin'
env GOOS=darwin GOARCH=amd64 go build -o binaries/darwin-amd64/discourse-reader
env GOOS=darwin GOARCH=arm64 go build -o binaries/darwin-arm64/discourse-reader

# compile 'dragonfly'
env GOOS=dragonfly GOARCH=amd64 go build -o binaries/dragonfly-amd64/discourse-reader

# compile 'freebsd'
env GOOS=freebsd GOARCH=amd64 go build -o binaries/freebsd-amd64/discourse-reader
env GOOS=freebsd GOARCH=arm64 go build -o binaries/freebsd-arm64/discourse-reader

# compile 'illumos'
env GOOS=illumos GOARCH=amd64 go build -o binaries/illumos-amd64/discourse-reader

# compile 'linux'
env GOOS=linux GOARCH=amd64 go build -o binaries/linux-amd64/discourse-reader
env GOOS=linux GOARCH=arm64 go build -o binaries/linux-arm64/discourse-reader
env GOOS=linux GOARCH=mips64 go build -o binaries/linux-mips64/discourse-reader
env GOOS=linux GOARCH=mips64le go build -o binaries/linux-mips64le/discourse-reader
env GOOS=linux GOARCH=ppc64 go build -o binaries/linux-ppc64/discourse-reader
env GOOS=linux GOARCH=ppc64le go build -o binaries/linux-ppc64le/discourse-reader
env GOOS=linux GOARCH=riscv64 go build -o binaries/linux-riscv64/discourse-reader
env GOOS=linux GOARCH=s390x go build -o binaries/linux-s390x/discourse-reader

# compile 'netbsd'
env GOOS=netbsd GOARCH=amd64 go build -o binaries/netbsd-amd64/discourse-reader
env GOOS=netbsd GOARCH=arm64 go build -o binaries/netbsd-arm64/discourse-reader

# compile 'openbsd'
env GOOS=openbsd GOARCH=amd64 go build -o binaries/openbsd-amd64/discourse-reader
env GOOS=openbsd GOARCH=arm64 go build -o binaries/openbsd-arm64/discourse-reader
env GOOS=openbsd GOARCH=mips64 go build -o binaries/openbsd-mips64/discourse-reader

# compile 'solaris'
env GOOS=solaris GOARCH=amd64 go build -o binaries/solaris-amd64/discourse-reader

# compile 'windows'
env GOOS=windows GOARCH=amd64 go build -o binaries/windows-amd64/discourse-reader.exe
env GOOS=windows GOARCH=386 go build -o binaries/windows-386/discourse-reader.exe
env GOOS=windows GOARCH=arm go build -o binaries/windows-arm/discourse-reader.exe
