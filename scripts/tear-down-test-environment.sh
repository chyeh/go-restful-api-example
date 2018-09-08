#!/bin/bash

./drop-db-schema.sh postgres://hellofresh:hellofresh@localhost:5432/hellofresh >/dev/null
docker rm -f $(docker ps -aqf "name=$1") >/dev/null