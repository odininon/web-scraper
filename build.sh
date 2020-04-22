#!/bin/bash

CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o scraper *.go