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

DATE=20160511

BETA=100
BEN=8

MUT=0.00001

COUNT=1
DIRNAME=$DATE-$COUNT
time $BIN/runsim8.sh $DIRNAME 414 -g $GENS -b $BEN -beta $BETA -singdef -passmutall -passmut $MUT
tar czf $DIRNAME.tar.gz $DIRNAME

COUNT=0
for MUT in 0.0001 0.001 0.01 0.1;
do
  COUNT=$((COUNT+1))
  DIRNAME=$DATE-$COUNT
  time $BIN/runsim8.sh $DIRNAME $SIMS -g $GENS -b $BEN -beta $BETA -singdef -passmutall -passmut $MUT
  tar czf $DIRNAME.tar.gz $DIRNAME
done
