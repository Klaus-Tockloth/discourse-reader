#!/bin/sh

# ------------------------------------
# Purpose:
# - Builds assets (tar.gz or zip).
#
# Releases:
# - v1.0.0 - 2022-11-04: initial release
# ------------------------------------

# set -o xtrace
set -o verbose

# recreate directory
rm -r ./assets
mkdir ./assets

# asset 'aix'
tar -cvzf ./assets/aix-ppc64_discourse-reader.tar.gz ./binaries/aix-ppc64/discourse-reader

# assets 'darwin'
tar -cvzf ./assets/darwin-amd64_discourse-reader.tar.gz ./binaries/darwin-amd64/discourse-reader
tar -cvzf ./assets/darwin-arm64_discourse-reader.tar.gz ./binaries/darwin-arm64/discourse-reader

# assets 'dragonfly'
tar -cvzf ./assets/dragonfly-amd64_discourse-reader.tar.gz ./binaries/dragonfly-amd64/discourse-reader

# assets 'freebsd'
tar -cvzf ./assets/freebsd-amd64_discourse-reader.tar.gz/freebsd-amd64/discourse-reader
tar -cvzf ./assets/freebsd-arm64_discourse-reader.tar.gz ./binaries/freebsd-arm64/discourse-reader

# asset 'illumos'
tar -cvzf ./assets/illumos-amd64_discourse-reader.tar.gz ./binaries/illumos-amd64/discourse-reader

# assets 'linux'
tar -cvzf ./assets/linux-amd64_discourse-reader.tar.gz ./binaries/linux-amd64/discourse-reader
tar -cvzf ./assets/linux-arm64_discourse-reader.tar.gz ./binaries/linux-arm64/discourse-reader
tar -cvzf ./assets/linux-mips64_discourse-reader.tar.gz ./binaries/linux-mips64/discourse-reader
tar -cvzf ./assets/linux-mips64le_discourse-reader.tar.gz ./binaries/linux-mips64le/discourse-reader
tar -cvzf ./assets/linux-ppc64_discourse-reader.tar.gz ./binaries/linux-ppc64/discourse-reader
tar -cvzf ./assets/linux-ppc64le_discourse-reader.tar.gz ./binaries/linux-ppc64le/discourse-reader
tar -cvzf ./assets/linux-riscv64_discourse-reader.tar.gz ./binaries/linux-riscv64/discourse-reader
tar -cvzf ./assets/linux-s390x_discourse-reader.tar.gz ./binaries/linux-s390x/discourse-reader

# assets 'netbsd'
tar -cvzf ./assets/netbsd-amd64_discourse-reader.tar.gz ./binaries/netbsd-amd64/discourse-reader
tar -cvzf ./assets/netbsd-arm64_discourse-reader.tar.gz ./binaries/netbsd-arm64/discourse-reader

# assets 'openbsd'
tar -cvzf ./assets/openbsd-amd64_discourse-reader.tar.gz ./binaries/openbsd-amd64/discourse-reader
tar -cvzf ./assets/openbsd-arm64_discourse-reader.tar.gz ./binaries/openbsd-arm64/discourse-reader
tar -cvzf ./assets/openbsd-mips64_discourse-reader.tar.gz ./binaries/openbsd-mips64/discourse-reader

# asset 'solaris'
tar -cvzf ./assets/solaris-amd64_discourse-reader.tar.gz ./binaries/solaris-amd64/discourse-reader

# assets 'windows'
zip ./assets/windows-amd64_discourse-reader.zip ./binaries/windows-amd64/discourse-reader.exe
zip ./assets/windows-386_discourse-reader.zip ./binaries/windows-386/discourse-reader.exe
zip ./assets/windows-arm_discourse-reader.zip ./binaries/windows-arm/discourse-reader.exe
