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

RESULTSDIR=$1
ROUNDS=$2

# create directory for results
mkdir $RESULTSDIR
cd ./$RESULTSDIR

# evaluate L8 BBB (aka, L8 Judging)
$BIN/runl8eval.sh pgg.L8JudgingAgent $ROUNDS

# evaluate L8 BBG
$BIN/runl8eval.sh pgg.L8BBGAgent $ROUNDS

# evaluate L8 BGB
$BIN/runl8eval.sh pgg.L8BGBAgent $ROUNDS

# evaluate L8 BGG
$BIN/runl8eval.sh pgg.L8BGGAgent $ROUNDS

# evaluate L8 GBB
$BIN/runl8eval.sh pgg.L8GBBAgent $ROUNDS

# evaluate L8 GBG
$BIN/runl8eval.sh pgg.L8GBGAgent $ROUNDS

# evaluate L8 GGB (aka, L8 Standing)
$BIN/runl8eval.sh pgg.L8StandingAgent $ROUNDS

# evaluate L8 GGG
$BIN/runl8eval.sh pgg.L8GGGAgent $ROUNDS

# done
cd ..
