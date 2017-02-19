#!/usr/bin/env bash

rm -rf drafthouse-seat-finder
git pull
CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" .
docker build -t jroyal/drafthouse-seat-finder .
docker stop drafthouse-seat-finder
docker container prune -f
docker image prune -f
docker run -d --name drafthouse-seat-finder -p 8080:8080 jroyal/drafthouse-seat-finder
