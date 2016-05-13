#!/bin/sh

#  run-gpgg-all.sh
#  
#
#  Created by John Maloney on 5/13/16.
#

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

# echo $BIN

# generate results for each graph type:
#   0: regular ring
#   1: homogeneous random
#   2: heterogeneous random (Watts-Strogatz w/ p=1.0)
#   3: scale free (with preferential attachment)
#   4: scale free (with uniform atachment)
#   5: small world (Watts-Strogatz w/ p=0.1)
#   6: small world (Watts-Strogatz w/ p=0.4)

for GTYPE in 2 3 4 0 1 5 6;
do
  DIR=gtype$GTYPE
  time $BIN/rungtype.sh $DIR $GTYPE 2 5
done