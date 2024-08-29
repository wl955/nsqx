#!/bin/bash

if [ -z "$1" ]; then
  echo "未提供topic"
  exit 1
fi

if [ -z "$2" ]; then
  echo "未提供channel"
  exit 1
fi

#新建一个topic
curl -X POST http://127.0.0.1:4151/topic/create?topic=$1

#新建一个channel
curl -X POST http://127.0.0.1:4151/channel/create?topic=$1\&channel=$2

echo 'ok'

exit 0
