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
GO=$BIN

# run the Python script to execute simulation
# assumption: the Python program is in the $PYT directory
BASE='stats'-$1
mkdir $BASE
cd $BASE
for i in `seq -f "%03g" 1 $2`;
do
  echo "  executing run $i..."
  $GO/runsim -g $3 -t 64 -a 64 -b 17 -c 1 -beta 10000 -mp -f $i.csv &> ./$i.log
done
cd ..
# python $PYT/runsim.py "$@"
