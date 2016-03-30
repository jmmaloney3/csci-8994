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

# iterate over values for Beta

# beta = 10^5
#BETA=100000
#time $BIN/runsim8.sh beta-$BETA $SIMS -g $GENS -beta $BETA

# beta = 10^3
#BETA=1000
#time $BIN/runsim8.sh beta-$BETA $SIMS -g $GENS -beta $BETA

# beta = 10^1
#BETA=10
#time $BIN/runsim8.sh beta-$BETA $SIMS -g $GENS -beta $BETA

BETA = 1000

# benefit = 5
BEN=5
time $BIN/runsim8.sh ben-$BEN $SIMS -g $GENS -b $BEN

# benefit = 10
BEN=10
time $BIN/runsim8.sh ben-$BEN $SIMS -g $GENS -b $BEN

# benefit = 20
BEN=20
time $BIN/runsim8.sh ben-$BEN $SIMS -g $GENS -b $BEN
	