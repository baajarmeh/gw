#!/bin/bash
#
# gw framework command interface
#

readonly workDir=$(pwd)
readonly gwSrcDir="$GOPATH/src/github.com/oceanho/gw"
readonly projectTemplate="$gwSrcDir/templates/scaffold/projectTemplate"
readonly applicationTemplate="$gwSrcDir/templates/scaffold/applicationTemplate"


function help() {
    echo "|-------------------------------------------------------------------------|"
    echo "| gw framework cli tools                                                  |"
    echo "|-------------------------------------------------------------------------|"
    echo "|                                                                         |"
    echo "| Usage:                                                                  |"
    echo "|    gwcli <Command> <options>                                            |"
    echo "| Commands:                                                               |"
    echo "|    newproject <Project Name> <Directory>, Create a gw project scaffold  |"
    echo "|    createapp <Application Name> <Directory>, Create a gw app scaffold   |"
    echo "|                                                                         |"
    echo "|--------------- gw cli v0.1  --------------------------------------------|"
}

function createProject() {
  local dir="$2"
  [ -z "$dir" ] && dir="$workDir"
  [ "$dir" = "." ] && dir="$workDir"
  local projectName="$1"
  if [ -z "$projectName" ]; then
    echo "Project Name are empty."
  else
    [ -d "$dir" ] || mkdir -p "$dir"
    if [ "`ls -A $dir`" = "" ]; then
      rsync -avrp $projectTemplate/ $dir 1>/dev/null
      for f in `find $dir -type f` ; do
        sed -i.bak -e "s#projectTemplate#$projectName#g" $f
        rm -f ${f}.bak
      done
      sleep 1
      local salt=$(openssl rand -base64 32)
      salt=${salt:0:32}
      local secret=$(openssl rand -base64 32)
      secret=${secret:0:32}
      sed -i.bak \
        -e "s#\$security-crypto-hash-salt#$salt#g" \
        -e "s#\$security-crypto-protect-secret#$secret#g" $dir/config/app.yaml
      rm -f $dir/config/app.yaml.bak
      cd "$dir" && \
      go mod init $projectName && go mod tidy && \
      git init "$dir" 1>/dev/null && git add ./ 1>/dev/null && git commit -m "Initial" 1>/dev/null && cd - 1>/dev/null && \
      echo "Create project scaffold into $dir Successfully."
    else
      echo "[WARN] $dir are not empty"
    fi
  fi
}

function createApp() {
  echo "create sub app"
}

function __execCommand() {
  local cmd="$1"
  case "$cmd" in
  createproject|newproject|startproject|project)
    shift 1
    createProject "$@"
    ;;
  createapp|newapp|startapp|app)
    shift 1
    createApp "$@"
    ;;
  *)
    help
    ;;
  esac
}

function main() {
  if [ $# -eq 0 ]; then
    help
  else
    __execCommand "$@"
  fi
}

main "$@"
