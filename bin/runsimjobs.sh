#!/bin/bash

# find the directory that holds the script
# - see http://stackoverflow.com/questions/59895/can-a-bash-script-tell-what-directory-its-stored-in/246128#246128
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
SOURCE="$(readlink "$SOURCE")"
[[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
# assumption: this bash script is in the $BIN directory
BIN="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# run several sim jobs

GENS=10000
#GENS=10
SIMS=500
#SIMS=5

DATE=20160420

BETA=5
BEN=5

DIRNAME=$DATE-1
time $BIN/runsim8.sh $DIRNAME $SIMS -g $GENS -b $BEN -beta $BETA -passmutall -am
tar czf $DIRNAME.tar.gz $DIRNAME

DIRNAME=$DATE-2
time $BIN/runsim8.sh $DIRNAME $SIMS -g $GENS -b $BEN -beta $BETA -singdef -am
tar czf $DIRNAME.tar.gz $DIRNAME

DIRNAME=$DATE-3
time $BIN/runsim8.sh $DIRNAME $SIMS -g $GENS -b $BEN -beta $BETA -am
tar czf $DIRNAME.tar.gz $DIRNAME

DIRNAME=$DATE-4
BETA=10
time $BIN/runsim8.sh $DIRNAME $SIMS -g $GENS -b $BEN -beta $BETA
tar czf $DIRNAME.tar.gz $DIRNAME
