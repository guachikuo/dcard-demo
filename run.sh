#!/bin/bash

docker-compose -p dcard-demo up -d --build --remove-orphans
docker rmi `docker images | grep "<none>" | awk {'print $3'}`