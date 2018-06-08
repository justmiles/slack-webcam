#!/bin/bash
source $GOPATH/src/gocv.io/x/gocv/env.sh
rm -rf build
mkdir -p build
go build main.go
mv main build/slack-webcam.linux-amd64
