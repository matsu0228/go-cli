#!/bin/sh

# Environment list
# $GOOS     $GOARCH
# darwin    386
# darwin    amd64
# freebsd   386
# freebsd   amd64
# freebsd   arm
# linux     386
# linux     amd64
# linux     arm
# netbsd    386
# netbsd    amd64
# netbsd    arm
# openbsd   386
# openbsd   amd64
# plan9     386
# plan9     amd64
# windows   386
# windows   amd64

set -e

GODIR="/home/dev/.gvm/gos/go1.9.1"
OS=("windows" "windows")
# OS=("darwin" "darwin" "freebsd" "freebsd" "freebsd" "linux" \
#   "linux" "linux" "netbsd" "netbsd" "netbsd" "openbsd" "openbsd" \
#   "plan9" "plan9" "windows" "windows")
ARCH=("386" "amd64") 
# ARCH=("386" "amd64" "386" "amd64" "arm" "386" "amd64" "arm" \
#   "386" "amd64" "arm" "386" "amd64" "386" "amd64" "386" "amd64")

#cd $(go env GOROOT)/src

cd ${GODIR}/src
pwd
go version
for i in `seq 0 1 16`
do
  GOOS=${OS["$i"]}
  GOARCH=${ARCH["$i"]}
  echo "\033[0;32m""[build-crosscompile-env.sh" ${GOOS} ${GOARCH} "] Build environment""\033[0m"
  GOOS=${GOOS} GOARCH=${GOARCH} ./make.bash
  echo "\033[0;32m""[build-crosscompile-env.sh" ${GOOS} ${GOARCH} "] done!""\033[0m"
done