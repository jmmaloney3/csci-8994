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

# remove the runsim command
EXE=$BIN/runsim
if [ -e $EXE ];
then
  echo "clean $EXE"
  rm $EXE
fi

# remove the rungpggsim command
EXE=$BIN/runsimgpgg
if [ -e $EXE ];
then
  echo "clean $EXE"
  rm $EXE
fi

# remove the sim package
PKG=$BIN/../go/pkg
if [ -e $PKG ];
then
  echo "clean $PKG"
  rm -r $PKG
fi
