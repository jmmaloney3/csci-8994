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

# evaluate punishment with voluntaty participation
# $BIN/runl8eval.sh pgg.Punisher $ROUNDS pgg.NonParticipant

# evaluate punishment with compulsory participation
$BIN/runl8eval.sh pgg.Punisher $ROUNDS

# evaluate voluntary participation without punishment
$BIN/runl8eval.sh pgg.NonParticipant $ROUNDS

# evaluate compulsory participation without punishment
# $BIN/runl8eval.sh pgg.L8BGGAgent $ROUNDS

# done
echo "DONE"
cd ..
