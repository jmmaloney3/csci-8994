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

# rename the files in directory $1 by appending $3 to the base filename
# Rename all *.csv to *b.csv
for f in $1/*.csv; do
mv -- "$f" "${f%.csv}b.csv"
done
mv -- $1/*.csv $2/

# Rename all *.log to *b.log
for f in $1/*.log; do
mv -- "$f" "${f%.log}b.log"
done
mv -- $1/*.log $2/
