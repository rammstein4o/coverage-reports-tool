#!/bin/sh

mkdir -p ./dist
docker build -t coverage-reports-tool:latest .
CONTAINER_ID=$(docker create coverage-reports-tool:latest)
docker cp ${CONTAINER_ID}:/coverage-reports-tool ./dist/
docker rm -v ${CONTAINER_ID}
