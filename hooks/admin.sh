#!/usr/bin/env bash

case $# in
   0)
      echo "Usage: $0 {endpoints|upgrade}"
      exit 1
      ;;
   1)
      case $1 in
         endpoints)
            HOOKS="$(cat /opt/webhook-linux-amd64/hooks.json | grep -o 'id.*' | cut -f2- -d:)"
            HOOKS="$(echo ${HOOKS::-1})"
            echo "[${HOOKS}]"
            ;;
         upgrade)
            rm -rf go/src/github.com/eschizoid/flixctl
            go get -u github.com/eschizoid/flixctl
            ;;
         *)
            echo "'$1' is not a valid admin command."
            echo "Usage: $0 {endpoints|upgrade}"
            exit 2
            ;;
      esac
      ;;
   *)
      echo "Usage: $0 {endpoints|upgrade}"
      exit 3
      ;;
esac
