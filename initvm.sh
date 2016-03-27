#!/bin/bash

sudo apt-get install -y git golang lzma
git clone https://github.com/jmmaloney3/csci-8994.git repo
cd repo
export GOPATH=$PWD/go/
go install sim
cd bin/
go build ../go/src/runsim.go
