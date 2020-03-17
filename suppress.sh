#!/bin/bash

if [ -z "$1" ]
  then
    echo "sleep time is required first arg"
    exit
fi

if [ -z "$2" ]
  then
    echo "url is required second arg"
    exit
fi

while :
do
    curl --request POST \
    --url http://localhost:3000/attack \
    --header 'content-type: application/json' \
    --data \
    "{
        \"url\":\"$2\",
        \"method\": \"GET\"
    }"

    echo "attacking every $1 seconds... ctrl+c to stop"
    sleep $1
done