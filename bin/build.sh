#!/bin/bash

# find the directory that holds the script
# - see http://stackoverflow.com/questions/59895/can-a-bash-script-tell-what-directory-its-stored-in/246128#246128
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
SOURCE="$(readlink "$SOURCE")"
[[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
BIN="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# echo $BIN

# assumption: this bash script is in the $BIN directory
# set the go path
export GOPATH=$BIN/../go

# build and install the sim package
go install sim
go install simgpgg

# run the tests
go test sim
go test simgpgg

# build the runsim command and put it in the bin directory
go build -o $BIN/runsim $GOPATH/src/runsim.go
go build -o $BIN/runsimgpgg $GOPATH/src/runsimgpgg.go
