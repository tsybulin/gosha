#!/bin/sh

env GOOS=linux GOARCH=amd64 go build -a -o . cmd/gosha.go

cd cmd
tar -cf ../config.tar config
cd ..

scp gosha config.tar pico:docker/gosha

rm gosha config.tar
