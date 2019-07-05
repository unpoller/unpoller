#!/bin/bash
# Usage:
# ./convert.sh <file> [prefix]
# <file> should contain a go struct, like uap_type.go
# It converts the go struct to an influx thing, like you see in uap_influx.go.
# [prefix] is optional. I used it to do all the stat_ uap metrics.
# Very crude, just helps skip a lot of copy/paste.
path=$1
pre=$2
while IFS='' read -r line; do
  name=$(echo "${line}" | awk '{print $1}')
  type=$(echo "${line}" | awk '{print $2}')
  json=$(echo "${line}" | awk '{print $3}')
  json=$(echo "${json}" | cut -d\" -f2)

  if [ "$json" != "" ] && [ "$name" != "" ]; then
    [[ "$type" = Flex* ]] && suf=.Val
    echo "\"${pre}${json}\": u.Stat.${name}${suf},"
  fi
done < ${path}
