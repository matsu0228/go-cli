#!/bin/sh

function usage_exit() {
  echo -e "
  Usage: $0 [-t TARGET] [-o OUTPUT] 
      GOのプログラムを他OS向けにコンパイル
      -t [target program]: デフォルト=./main.go
      -o [output dir/name]: デフォルト=./main**
  " 1>&2
  exit 1
}

function set_color(){
  echo -e "$1"
}
COLOR_ALERT='\033[0;31m'
COLOR_INFO='\033[32m'
COLOR_END='\033[0m'

while getopts to:h OPT
do
  case $OPT in
    t)  OPT_TARGET_FILE=${OPTARG}
      ;;
    o)  OPT_OUTPUT_FILE=${OPTARG}
      ;;
    *) usage_exit
      ;;
  esac
done
shift $((OPTIND - 1))

# validate OPTS
if [ -z "${OPT_TARGET_FILE}" ]; then
  OPT_TARGET_FILE="./main.go"
  set_color ${COLOR_INFO}
  echo -e "ターゲットの指定がないため、./main.goをbuildします"
  set_color ${COLOR_END}
fi
if [ -z "${OPT_OUTPUT_FILE}" ]; then
  OUTPUT_BASE="./main"
else
  OUTPUT_BASE="${OPT_OUTPUT_FILE}"
fi

# end of perser --------------------------------

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


GODIR="/home/dev/.gvm/gos/go1.9.1"
OS=("windows" "windows")
# OS=("darwin" "darwin" "freebsd" "freebsd" "freebsd" "linux" \
#   "linux" "linux" "netbsd" "netbsd" "netbsd" "openbsd" "openbsd" \
#   "plan9" "plan9" "windows" "windows")
ARCH=("386" "amd64")
# ARCH=("386" "amd64" "386" "amd64" "arm" "386" "amd64" "arm" \
#   "386" "amd64" "arm" "386" "amd64" "386" "amd64" "386" "amd64")
EXTENT=(".exe" ".exe") #出力ファイルの拡張子 win:.exe / linux:nil

#cd $(go env GOROOT)/src

UPPER_COMPILE=0
if [ ${#OS[@]} -eq ${#ARCH[@]} ]; then
  UPPER_COMPILE=$((${#OS[@]}-1))
else
  usage_exit
fi

go version
echo "compile target is ${OPT_TARGET_FILE}"

set_color ${COLOR_INFO}
for i in `seq 0 1 ${UPPER_COMPILE}`
do
  echo "loop of $i"
  GOOS=${OS["$i"]}
  GOARCH=${ARCH["$i"]}
  GO_BUILD_OPT=" -o ${OUTPUT_BASE}_${GOOS}_${GOARCH}${EXTENT["$i"]}"
  echo "[compile: " ${GOOS} ${GOARCH} "] Build environment with ${GO_BUILD_OPT}"
  env GOOS=${GOOS} GOARCH=${GOARCH} go build ${GO_BUILD_OPT}
done
set_color ${COLOR_END}